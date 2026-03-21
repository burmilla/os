package app

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/burmilla/os/pkg/libcompose/cli/logger"
	"github.com/burmilla/os/pkg/libcompose/docker"
	"github.com/burmilla/os/pkg/libcompose/project"
)

// ProjectFactory is a struct that holds the app.ProjectFactory implementation.
type ProjectFactory struct {
}

// Create implements ProjectFactory.Create using docker client.
func (p *ProjectFactory) Create(c *cli.Context) (project.APIProject, error) {
	context := &docker.Context{}
	context.LoggerFactory = logger.NewColorLoggerFactory()
	Populate(context, c)

	context.ComposeFiles = c.GlobalStringSlice("file")

	if len(context.ComposeFiles) == 0 {
		context.ComposeFiles = []string{"docker-compose.yml"}
		if _, err := os.Stat("docker-compose.override.yml"); err == nil {
			context.ComposeFiles = append(context.ComposeFiles, "docker-compose.override.yml")
		}
	}

	context.ProjectName = c.GlobalString("project-name")

	return docker.NewProject(context, nil)
}
