package docker

import (
	"github.com/burmilla/os/pkg/util"

	composeConfig "github.com/burmilla/os/pkg/libcompose/config"
	"github.com/burmilla/os/pkg/libcompose/docker"
	"github.com/burmilla/os/pkg/libcompose/project"
)

type ServiceFactory struct {
	Context *docker.Context
	Deps    map[string][]string
}

// depAliases maps legacy service names to their actual service names.
// This allows cloud-config authors to use well-known names in
// io.rancher.os.after / io.rancher.os.before labels without needing
// to know internal implementation details.
var depAliases = map[string]string{
	"cloud-init": "cloud-init-execute",
}

func (s *ServiceFactory) Create(project *project.Project, name string, serviceConfig *composeConfig.ServiceConfig) (project.Service, error) {
	labels := serviceConfig.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	if after := labels["io.rancher.os.after"]; after != "" {
		for _, dep := range util.TrimSplit(after, ",") {
			if alias, ok := depAliases[dep]; ok {
				dep = alias
			}
			s.Deps[name] = append(s.Deps[name], dep)
		}
	}
	if before := labels["io.rancher.os.before"]; before != "" {
		for _, dep := range util.TrimSplit(before, ",") {
			if alias, ok := depAliases[dep]; ok {
				dep = alias
			}
			s.Deps[dep] = append(s.Deps[dep], name)
		}
	}

	return NewService(s, name, serviceConfig, s.Context, project), nil
}
