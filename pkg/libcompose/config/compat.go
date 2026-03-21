package config

import (
	"fmt"
	"strings"

	composetypes "github.com/compose-spec/compose-go/types"
)

// ServiceConfigV1ToServiceConfig converts a v1 service config to a compose-go ServiceConfig.
func ServiceConfigV1ToServiceConfig(v1 *ServiceConfigV1) *composetypes.ServiceConfig {
	if v1 == nil {
		return nil
	}

	sc := &composetypes.ServiceConfig{
		CapAdd:        v1.CapAdd,
		CapDrop:       v1.CapDrop,
		CgroupParent:  v1.CgroupParent,
		CPUQuota:      v1.CPUQuota,
		CPUSet:        v1.CPUSet,
		CPUShares:     v1.CPUShares,
		Command:       composetypes.ShellCommand(v1.Command),
		ContainerName: v1.ContainerName,
		Devices:       v1.Devices,
		DNS:           composetypes.StringList(v1.DNS),
		DNSSearch:     composetypes.StringList(v1.DNSSearch),
		DomainName:    v1.DomainName,
		Entrypoint:    composetypes.ShellCommand(v1.Entrypoint),
		EnvFile:       composetypes.StringList(v1.EnvFile),
		Environment:   sliceToMappingWithEquals(v1.Environment),
		Expose:        stringSliceToStringOrNumberList(v1.Expose),
		ExternalLinks: v1.ExternalLinks,
		ExtraHosts:    sliceToHostsList(v1.ExtraHosts),
		Hostname:      v1.Hostname,
		Image:         v1.Image,
		Labels:        composetypes.Labels(v1.Labels),
		Links:         v1.Links,
		LogDriver:     v1.LogDriver,
		LogOpt:        v1.LogOpt,
		MacAddress:    v1.MacAddress,
		MemLimit:      composetypes.UnitBytes(v1.MemLimit),
		MemSwapLimit:  composetypes.UnitBytes(v1.MemSwapLimit),
		Net:           v1.Net,
		NetworkMode:   v1.Net,
		OomScoreAdj:   int64(v1.OomScoreAdj),
		Pid:           v1.Pid,
		Ports:         stringSliceToPortConfigs(v1.Ports),
		Privileged:    v1.Privileged,
		Restart:       v1.Restart,
		ReadOnly:      v1.ReadOnly,
		StdinOpen:     v1.StdinOpen,
		SecurityOpt:   v1.SecurityOpt,
		Tty:           v1.Tty,
		User:          v1.User,
		Uts:           v1.Uts,
		Ipc:           v1.Ipc,
		VolumeDriver:  v1.VolumeDriver,
		Volumes:       stringSliceToVolumeConfigs(v1.Volumes),
		VolumesFrom:   v1.VolumesFrom,
		WorkingDir:    v1.WorkingDir,
	}

	if v1.LogDriver != "" || v1.LogOpt != nil {
		sc.Logging = &composetypes.LoggingConfig{
			Driver:  v1.LogDriver,
			Options: v1.LogOpt,
		}
	}

	if v1.Build != "" {
		sc.Build = &composetypes.BuildConfig{
			Context:    v1.Build,
			Dockerfile: v1.Dockerfile,
		}
	}

	if v1.Ulimits.Elements != nil {
		sc.Ulimits = make(map[string]*composetypes.UlimitsConfig)
		for _, u := range v1.Ulimits.Elements {
			sc.Ulimits[u.Name] = &composetypes.UlimitsConfig{
				Soft: int(u.Soft),
				Hard: int(u.Hard),
			}
		}
	}

	return sc
}

// ConvertV1Services converts a map of v1 service configs to compose-go ServiceConfigs.
func ConvertV1Services(v1Services map[string]*ServiceConfigV1) map[string]*composetypes.ServiceConfig {
	result := make(map[string]*composetypes.ServiceConfig, len(v1Services))
	for name, v1 := range v1Services {
		sc := ServiceConfigV1ToServiceConfig(v1)
		sc.Name = name
		result[name] = sc
	}
	return result
}

// VolumesToStringSlice converts compose-go ServiceVolumeConfig slice to simple strings.
func VolumesToStringSlice(volumes []composetypes.ServiceVolumeConfig) []string {
	result := make([]string, 0, len(volumes))
	for _, v := range volumes {
		result = append(result, volumeConfigToString(v))
	}
	return result
}

// PortsToStringSlice converts compose-go ServicePortConfig slice to simple port strings.
func PortsToStringSlice(ports []composetypes.ServicePortConfig) []string {
	result := make([]string, 0, len(ports))
	for _, p := range ports {
		result = append(result, portConfigToString(p))
	}
	return result
}

// EnvironmentToSlice converts MappingWithEquals to KEY=VALUE string slice.
func EnvironmentToSlice(env composetypes.MappingWithEquals) []string {
	result := make([]string, 0, len(env))
	for k, v := range env {
		if v != nil {
			result = append(result, fmt.Sprintf("%s=%s", k, *v))
		} else {
			result = append(result, k)
		}
	}
	return result
}

// LoggingDriver returns the logging driver from a ServiceConfig, handling nil.
func LoggingDriver(sc *composetypes.ServiceConfig) string {
	if sc.Logging != nil {
		return sc.Logging.Driver
	}
	return sc.LogDriver
}

// LoggingOptions returns the logging options from a ServiceConfig, handling nil.
func LoggingOptions(sc *composetypes.ServiceConfig) map[string]string {
	if sc.Logging != nil {
		return sc.Logging.Options
	}
	return sc.LogOpt
}

func sliceToMappingWithEquals(slice []string) composetypes.MappingWithEquals {
	if slice == nil {
		return nil
	}
	return composetypes.NewMappingWithEquals(slice)
}

func stringSliceToStringOrNumberList(slice []string) composetypes.StringOrNumberList {
	result := make(composetypes.StringOrNumberList, len(slice))
	for i, s := range slice {
		result[i] = s
	}
	return result
}

func sliceToHostsList(slice []string) composetypes.HostsList {
	if slice == nil {
		return nil
	}
	hosts := composetypes.HostsList{}
	for _, s := range slice {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) == 2 {
			hosts[parts[0]] = parts[1]
		}
	}
	return hosts
}

func stringSliceToPortConfigs(ports []string) []composetypes.ServicePortConfig {
	result := make([]composetypes.ServicePortConfig, 0, len(ports))
	for _, p := range ports {
		result = append(result, parsePortString(p))
	}
	return result
}

func parsePortString(s string) composetypes.ServicePortConfig {
	// Simple parse: handle "hostPort:containerPort" and "containerPort"
	pc := composetypes.ServicePortConfig{
		Protocol: "tcp",
	}

	// Check for protocol suffix first
	proto := "tcp"
	portPart := s
	if idx := strings.LastIndex(s, "/"); idx >= 0 {
		proto = s[idx+1:]
		portPart = s[:idx]
	}
	pc.Protocol = proto

	parts := strings.SplitN(portPart, ":", 3)
	switch len(parts) {
	case 1:
		// "containerPort"
		fmt.Sscanf(parts[0], "%d", &pc.Target)
	case 2:
		// "hostPort:containerPort"
		pc.Published = parts[0]
		fmt.Sscanf(parts[1], "%d", &pc.Target)
	case 3:
		// "ip:hostPort:containerPort"
		pc.HostIP = parts[0]
		pc.Published = parts[1]
		fmt.Sscanf(parts[2], "%d", &pc.Target)
	}

	return pc
}

func portConfigToString(p composetypes.ServicePortConfig) string {
	var s string
	if p.HostIP != "" {
		s = fmt.Sprintf("%s:", p.HostIP)
	}
	if p.Published != "" {
		s += fmt.Sprintf("%s:", p.Published)
	}
	s += fmt.Sprintf("%d", p.Target)
	if p.Protocol != "" && p.Protocol != "tcp" {
		s += "/" + p.Protocol
	}
	return s
}

func stringSliceToVolumeConfigs(volumes []string) []composetypes.ServiceVolumeConfig {
	result := make([]composetypes.ServiceVolumeConfig, 0, len(volumes))
	for _, v := range volumes {
		result = append(result, parseVolumeString(v))
	}
	return result
}

func parseVolumeString(s string) composetypes.ServiceVolumeConfig {
	parts := strings.SplitN(s, ":", 3)
	vc := composetypes.ServiceVolumeConfig{
		Type: composetypes.VolumeTypeBind,
	}

	switch len(parts) {
	case 1:
		vc.Target = parts[0]
		if !strings.HasPrefix(parts[0], "/") && !strings.HasPrefix(parts[0], ".") {
			vc.Type = composetypes.VolumeTypeVolume
			vc.Source = parts[0]
		}
	case 2:
		vc.Source = parts[0]
		vc.Target = parts[1]
	case 3:
		vc.Source = parts[0]
		vc.Target = parts[1]
		if parts[2] == "ro" {
			vc.ReadOnly = true
		}
	}

	return vc
}

func volumeConfigToString(v composetypes.ServiceVolumeConfig) string {
	s := ""
	if v.Source != "" {
		s = v.Source + ":"
	}
	s += v.Target
	if v.ReadOnly {
		s += ":ro"
	}
	return s
}
