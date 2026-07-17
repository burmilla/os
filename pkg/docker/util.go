package docker

import (
	"github.com/burmilla/os/config"

	composeConfig "github.com/burmilla/os/pkg/libcompose/config"
)

func IsSystemContainer(serviceConfig *composeConfig.ServiceConfig) bool {
	return serviceConfig.Labels[config.ScopeLabel] == config.System
}
