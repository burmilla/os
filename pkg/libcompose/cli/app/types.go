package app

import (
	"github.com/burmilla/os/pkg/libcompose/project"
	"github.com/codegangsta/cli"
)

// ProjectFactory is an interface that helps creating libcompose project.
type ProjectFactory interface {
	// Create creates a libcompose project from the command line options (codegangsta cli context).
	Create(c *cli.Context) (project.APIProject, error)
}
