package docker

import (
	"github.com/burmilla/os/pkg/libcompose/project"
	"github.com/docker/docker/cliconfig"
)

// Context holds context meta information about a libcompose project and docker
// client information (like configuration file, builder to use, …)
type Context struct {
	project.Context
	ClientFactory project.ClientFactory
	ConfigDir     string
	ConfigFile    *cliconfig.ConfigFile
	AuthLookup    AuthLookup
}

func (c *Context) open() error {
	return c.LookupConfig()
}

// LookupConfig tries to load the docker configuration files, if any.
func (c *Context) LookupConfig() error {
	if c.ConfigFile != nil {
		return nil
	}

	config, err := cliconfig.Load(c.ConfigDir)
	if err != nil {
		return err
	}

	c.ConfigFile = config

	return nil
}
