package compose

import (
	"fmt"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/docker"
	"github.com/burmilla/os/pkg/log"
	"github.com/burmilla/os/pkg/util"
	"github.com/burmilla/os/pkg/util/network"
	"github.com/burmilla/os/pkg/util/versions"

	composeConfig "github.com/burmilla/os/pkg/libcompose/config"
	"github.com/burmilla/os/pkg/libcompose/project"
	yaml "github.com/cloudfoundry-incubator/candiedyaml"
)

func LoadService(p *project.Project, cfg *config.CloudConfig, useNetwork bool, service string) error {
	// First check the multi engine service file.
	// If the name has been found in multi enging service file and matches, will not execute network.LoadServiceResource
	// Otherwise will execute network.LoadServiceResource
	bytes, err := network.LoadMultiEngineResource(service)
	if err != nil || bytes == nil {
		bytes, err = network.LoadServiceResource(service, useNetwork, cfg)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	m := map[interface{}]interface{}{}
	if err = yaml.Unmarshal(bytes, &m); err != nil {
		e := fmt.Errorf("Failed to parse YAML configuration for %s: %v", service, err)
		log.Error(e)
		return e
	}

	if engineType := cfg.SelectedEngineType(); engineType != config.EngineDocker && userService(m) {
		e := fmt.Errorf("service %s runs in user Docker, which is not available with the %s engine", service, engineType)
		log.Error(e)
		return e
	}

	m = adjustContainerNames(m)

	bytes, err = yaml.Marshal(m)
	if err != nil {
		e := fmt.Errorf("Failed to marshal YAML configuration for %s: %v", service, err)
		log.Error(e)
		return e
	}

	if err = p.Load(bytes); err != nil {
		e := fmt.Errorf("Failed to load %s: %v", service, err)
		log.Error(e)
		return e
	}

	return nil
}

func LoadSpecialService(p *project.Project, cfg *config.CloudConfig, serviceName, serviceValue string) error {
	if serviceValue == "" {
		return fmt.Errorf("empty service value for %q", serviceName)
	}

	// Save config in case load fails so we can roll back
	previousConfig, ok := p.ServiceConfigs.Get(serviceName)

	p.ServiceConfigs.Add(serviceName, &composeConfig.ServiceConfig{})

	if err := LoadService(p, cfg, true, serviceValue); err != nil {
		// Rollback to previous config to avoid leaving an empty placeholder
		if ok {
			p.ServiceConfigs.Add(serviceName, previousConfig)
		}
		return err
	}

	return nil
}

func loadConsoleService(cfg *config.CloudConfig, p *project.Project) error {
	if cfg.Rancher.Console == "" || cfg.Rancher.Console == "default" {
		return nil
	}
	return LoadSpecialService(p, cfg, "console", cfg.Rancher.Console)
}

// LoadEngine loads the compose definition of an engine (a name from the
// "engines" list of the service repositories, a local file or an URL) into p
// and returns the name of the service it defines, which is the engine type:
// docker, k3s, ... The previous definition of that service, e.g. the built-in
// docker service, is replaced.
func LoadEngine(p *project.Project, cfg *config.CloudConfig, engine string) (string, error) {
	if engine == "" {
		return "", fmt.Errorf("empty engine name")
	}

	bytes, err := network.LoadServiceResource(engine, true, cfg)
	if err != nil {
		return "", err
	}

	m := map[interface{}]interface{}{}
	if err = yaml.Unmarshal(bytes, &m); err != nil {
		return "", fmt.Errorf("Failed to parse YAML configuration for engine %s: %v", engine, err)
	}

	serviceName, err := engineServiceName(engine, m)
	if err != nil {
		return "", err
	}

	m = adjustContainerNames(m)
	bytes, err = yaml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal YAML configuration for engine %s: %v", engine, err)
	}

	// Save config in case load fails so we can roll back
	previousConfig, ok := p.ServiceConfigs.Get(serviceName)
	p.ServiceConfigs.Add(serviceName, &composeConfig.ServiceConfig{})

	if err = p.Load(bytes); err != nil {
		if ok {
			p.ServiceConfigs.Add(serviceName, previousConfig)
		}
		return "", fmt.Errorf("Failed to load engine %s: %v", engine, err)
	}

	return serviceName, nil
}

// engineServiceName returns the single service an engine compose file defines.
func engineServiceName(engine string, m map[interface{}]interface{}) (string, error) {
	if len(m) != 1 {
		return "", fmt.Errorf("engine %s must define exactly one service, found %d", engine, len(m))
	}
	for k := range m {
		name, ok := k.(string)
		if !ok || name == "" {
			break
		}
		return name, nil
	}
	return "", fmt.Errorf("engine %s has an invalid service name", engine)
}

func loadEngineService(cfg *config.CloudConfig, p *project.Project) error {
	engine := cfg.SelectedEngine()

	// The configured engine may be missing from the service repository of
	// this OS version (e.g. an engine selected on an older release). Fall
	// back to another engine of the same type instead of leaving the system
	// without an engine. A bare type such as "docker" is not an engine in the
	// repository (there is no os-docker:latest image), so it always resolves
	// to the latest available engine of that type.
	engines, err := network.GetEngines(cfg.Rancher.Repositories.ToArray())
	if err == nil {
		resolved, exact := resolveEngine(engine, engines)
		if !exact && !config.IsEngineType(engine) {
			log.Warnf("Engine %s is not available in the service repositories, falling back to %s. Use 'ros engine list' and 'ros engine switch' to select an available engine.", engine, resolved)
		}
		engine = resolved
	}
	if config.IsEngineType(engine) {
		// Engine index could not be read or has no engine of this type, keep
		// the built-in service.
		return nil
	}

	_, err = LoadEngine(p, cfg, engine)
	return err
}

// resolveEngine returns the engine to load and whether it is the requested
// one. The requested engine is loaded as is when it is a local file or URL,
// when it is listed in available or when the engine index could not be read
// (available is empty). A bare engine type resolves to the latest available
// engine of that type. An engine missing from available falls back to the
// latest available engine of its type, else to the latest Docker engine, else
// to the latest engine of any type. A bare type without any engine of that
// type is returned unchanged.
func resolveEngine(engine string, available []string) (string, bool) {
	if len(available) == 0 || config.IsLocalOrURL(engine) || util.Contains(available, engine) {
		return engine, true
	}

	engineType := config.EngineType(engine)
	if latest := latestEngine(available, engineType); latest != "" {
		return latest, false
	}
	if config.IsEngineType(engine) {
		return engine, true
	}
	if latest := latestEngine(available, config.EngineDocker); latest != "" {
		return latest, false
	}
	return latestEngine(available, ""), false
}

// latestEngine returns the latest engine of the given type in available, or
// of any type when engineType is empty. It returns "" when there is none.
func latestEngine(available []string, engineType string) string {
	candidates := []string{}
	for _, engine := range available {
		if engineType == "" || config.EngineType(engine) == engineType {
			candidates = append(candidates, engine)
		}
	}
	if len(candidates) == 0 {
		return ""
	}
	versions.SortEngines(candidates)
	return candidates[len(candidates)-1]
}

// engineServices removes from services the built-in engine services of
// other types than the selected engine, e.g. the docker service when K3s is
// the selected engine, so that only one engine is started. Services that run
// in user Docker are removed too when the engine is not Docker.
func engineServices(cfg *config.CloudConfig, services map[string]*composeConfig.ServiceConfigV1) map[string]*composeConfig.ServiceConfigV1 {
	engineType := cfg.SelectedEngineType()
	result := map[string]*composeConfig.ServiceConfigV1{}
	for name, service := range services {
		if config.IsEngineType(name) && name != engineType {
			continue
		}
		if engineType != config.EngineDocker && service != nil && service.Labels[config.ScopeLabel] != config.System {
			log.Warnf("Not starting service %s: it runs in user Docker, which is not available with the %s engine", name, engineType)
			continue
		}
		result[name] = service
	}
	return result
}

// userService reports whether the service configuration m defines a service
// that runs in user Docker rather than in System Docker.
func userService(m map[interface{}]interface{}) bool {
	for _, v := range m {
		service, ok := v.(map[interface{}]interface{})
		if !ok {
			continue
		}
		labels, ok := service["labels"].(map[interface{}]interface{})
		if !ok || labels[config.ScopeLabel] != config.System {
			return true
		}
	}
	return false
}

func projectReload(p *project.Project, useNetwork *bool, loadConsole bool, environmentLookup *docker.ConfigEnvironment, authLookup *docker.ConfigAuthLookup) func() error {
	enabled := map[interface{}]interface{}{}
	return func() error {
		cfg := config.LoadConfig()

		environmentLookup.SetConfig(cfg)
		authLookup.SetConfig(cfg)

		enabled = addServices(p, enabled, engineServices(cfg, cfg.Rancher.Services))

		for service, serviceEnabled := range cfg.Rancher.ServicesInclude {
			if _, ok := enabled[service]; ok || !serviceEnabled {
				continue
			}

			if err := LoadService(p, cfg, *useNetwork, service); err != nil {
				if err != network.ErrNoNetwork {
					log.Errorf("Failed to load service(%s): %v", service, err)
				}
				continue
			}

			enabled[service] = service
		}

		if !*useNetwork {
			return nil
		}

		if loadConsole {
			if err := loadConsoleService(cfg, p); err != nil {
				log.Errorf("Failed to load rancher.console=(%s): %v", cfg.Rancher.Console, err)
			}
		}

		if err := loadEngineService(cfg, p); err != nil {
			log.Errorf("Failed to load engine %s: %v", cfg.SelectedEngine(), err)
		}

		return nil
	}
}
