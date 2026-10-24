package config

import "strings"

// Engine types. An engine is the runtime that runs the user's workloads on
// top of System Docker: user Docker by default. Engines are published in the
// service repositories (burmilla/os-services) under the "engines" key of
// index.yml as <type>-<version>, e.g. docker-29.9.2, and the compose file of
// an engine defines one service named after its type, e.g. "docker".
const (
	EngineDocker = "docker"
	EngineK3s    = "k3s"
	EngineQemu   = "qemu"
)

// EngineTypes lists the engine types ros knows how to start.
var EngineTypes = []string{EngineDocker, EngineK3s, EngineQemu}

// IsEngineType reports whether name is a bare engine type such as "docker",
// which selects the latest available engine of that type.
func IsEngineType(name string) bool {
	for _, t := range EngineTypes {
		if name == t {
			return true
		}
	}
	return false
}

// IsLocalOrURL reports whether service is a local file or an URL instead of a
// name in the service repositories.
func IsLocalOrURL(service string) bool {
	return strings.HasPrefix(service, "/") || strings.HasPrefix(service, "http:/") || strings.HasPrefix(service, "https:/")
}

// EngineType returns the type of an engine name: the part before the first
// "-" (docker-29.9.2 -> docker), or the name itself for a bare type. The
// legacy value "latest" and local files or URLs (custom Docker engines) are
// Docker.
func EngineType(engine string) string {
	if engine == "" || engine == "latest" || IsLocalOrURL(engine) {
		return EngineDocker
	}
	if i := strings.Index(engine, "-"); i > 0 {
		return engine[:i]
	}
	return engine
}

// SelectedEngine returns the engine selected in the configuration:
// rancher.engine, else the legacy rancher.docker.engine, else the defaults.
// "latest" is the legacy spelling of the latest Docker engine. The result is
// never empty: the built-in default is the latest Docker engine.
func (c *CloudConfig) SelectedEngine() string {
	for _, engine := range []string{
		c.Rancher.Engine,
		c.Rancher.Docker.Engine,
		c.Rancher.Defaults.Engine,
		c.Rancher.Defaults.Docker.Engine,
	} {
		if engine == "latest" {
			return EngineDocker
		}
		if engine != "" {
			return engine
		}
	}
	return EngineDocker
}

// SelectedEngineType returns the type of the selected engine.
func (c *CloudConfig) SelectedEngineType() string {
	return EngineType(c.SelectedEngine())
}
