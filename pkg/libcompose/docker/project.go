package docker

import (
	"os"
	"path/filepath"

	"github.com/burmilla/os/pkg/libcompose/config"
	"github.com/burmilla/os/pkg/libcompose/docker/client"
	"github.com/burmilla/os/pkg/libcompose/lookup"
	"github.com/burmilla/os/pkg/libcompose/project"
	"github.com/sirupsen/logrus"
)

// ComposeVersion is name of docker-compose.yml file syntax supported version
const ComposeVersion = "1.5.0"

// NewProject creates a Project with the specified context.
func NewProject(context *Context, parseOptions *config.ParseOptions) (*project.Project, error) {
	if context.ResourceLookup == nil {
		context.ResourceLookup = &lookup.FileConfigLookup{}
	}

	if context.EnvironmentLookup == nil {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		context.EnvironmentLookup = &lookup.ComposableEnvLookup{
			Lookups: []config.EnvironmentLookup{
				&lookup.EnvfileLookup{
					Path: filepath.Join(cwd, ".env"),
				},
				&lookup.OsEnvLookup{},
			},
		}
	}

	if context.AuthLookup == nil {
		context.AuthLookup = NewConfigAuthLookup(context)
	}

	if context.ServiceFactory == nil {
		context.ServiceFactory = &ServiceFactory{
			context: context,
		}
	}

	if context.ClientFactory == nil {
		factory, err := project.NewDefaultClientFactory(client.Options{})
		if err != nil {
			return nil, err
		}
		context.ClientFactory = factory
	}

	// FIXME(vdemeester) Remove the context duplication ?
	p := project.NewProject(context.ClientFactory, &context.Context, parseOptions)

	err := p.Parse()
	if err != nil {
		return nil, err
	}

	if err = context.open(); err != nil {
		logrus.Errorf("Failed to open project %s: %v", p.Name, err)
		return nil, err
	}

	return p, err
}
