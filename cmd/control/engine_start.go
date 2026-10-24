package control

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/log"

	"github.com/codegangsta/cli"
)

// engineStarter starts an engine of one type. start runs in the engine's
// System Docker container after its binaries were copied to
// /var/lib/rancher/engine, and normally execs 'ros engine-init <type>' in
// the console. init runs in the console once it is ready, and execs the
// engine daemon.
type engineStarter struct {
	start func(cfg *config.CloudConfig) error
	init  func(cfg *config.CloudConfig, args []string) error
	// links are the console paths to create for the engine binaries, in
	// addition to one /usr/bin link per file in /var/lib/rancher/engine.
	links []symlink
	// scripts returns wrapper scripts to write into the console, by path,
	// for engines whose tools run in the engine container.
	scripts func(cfg *config.CloudConfig) map[string]string
}

var engineStarters = map[string]engineStarter{
	// 'ros engine-start docker' is the same as the historical 'ros user-docker'
	config.EngineDocker: {
		start: startDocker,
		init:  func(_ *config.CloudConfig, args []string) error { return dockerInit(args) },
	},
}

// engineStartAction is 'ros engine-start <type>': the command of an engine's
// System Docker service.
func engineStartAction(c *cli.Context) error {
	starter, err := engineStarterFor(c.Args())
	if err != nil {
		return err
	}

	if err := prepareEngine(); err != nil {
		return err
	}

	return starter.start(config.LoadConfig())
}

// engineInitAction is 'ros engine-init <type> [args]', run inside the
// console by engine-start.
func engineInitAction(c *cli.Context) error {
	starter, err := engineStarterFor(c.Args())
	if err != nil {
		return err
	}

	waitForConsole()

	// Mounts made by the engine (overlays, volumes) must propagate to the
	// host and the other containers that share the console's filesystem.
	if err := syscall.Mount("", "/", "", syscall.MS_SHARED|syscall.MS_REC, ""); err != nil {
		log.Error(err)
	}
	if err := syscall.Mount("", "/run", "", syscall.MS_SHARED|syscall.MS_REC, ""); err != nil {
		log.Error(err)
	}

	cfg := config.LoadConfig()
	linkEngineBinaries(cfg, c.Args()[0])

	return starter.init(cfg, c.Args()[1:])
}

func engineStarterFor(args []string) (engineStarter, error) {
	if len(args) < 1 {
		return engineStarter{}, fmt.Errorf("Must specify the engine type")
	}
	starter, ok := engineStarters[args[0]]
	if !ok {
		return engineStarter{}, fmt.Errorf("%s is not a supported engine type", args[0])
	}
	return starter, nil
}

// startInConsole execs 'ros engine-init <type>' in the console, which is the
// usual start step of an engine.
func startInConsole(cfg *config.CloudConfig, engineType string, env []string) error {
	return execInConsole(cfg, "", append(os.Environ(), env...), []string{"ros", "engine-init", engineType})
}

func waitForConsole() {
	// TODO: this should be replaced by a "Console ready event watcher"
	for {
		if _, err := os.Stat(consoleDone); err == nil {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// engineBinarySymlinks returns the console links for the binaries of an
// engine type other than Docker (Docker has its own list in
// symLinkEngineBinary): one /usr/bin link per file in /var/lib/rancher/engine
// plus the links the engine type declares.
func engineBinarySymlinks(engineType string) []symlink {
	starter, ok := engineStarters[engineType]
	if !ok || engineType == config.EngineDocker {
		return nil
	}

	links := []symlink{}
	files, err := ioutil.ReadDir(destDirectory)
	if err != nil {
		log.Debugf("No engine binaries in %s: %v", destDirectory, err)
	}
	for _, file := range files {
		if file.IsDir() || file.Name() == dockerCompletionFName {
			continue
		}
		links = append(links, symlink{path.Join(destDirectory, file.Name()), path.Join("/usr/bin", file.Name())})
	}

	return append(links, starter.links...)
}

// linkEngineBinaries makes the tools of an engine type available in the
// console: links to its binaries and wrapper scripts.
func linkEngineBinaries(cfg *config.CloudConfig, engineType string) {
	for _, link := range engineBinarySymlinks(engineType) {
		syscall.Unlink(link.newname)
		if _, err := os.Stat(link.oldname); err != nil {
			continue
		}
		if err := os.Symlink(link.oldname, link.newname); err != nil {
			log.Error(err)
		}
	}
	for file, content := range engineScripts(cfg, engineType) {
		if err := ioutil.WriteFile(file, []byte(content), 0755); err != nil {
			log.Error(err)
		}
	}
}

// unlinkEngineBinaries removes what linkEngineBinaries created, when the
// engine type is switched. Links to binaries of the Docker engine are kept:
// they dangle until Docker runs again, like before Docker started.
func unlinkEngineBinaries(cfg *config.CloudConfig, engineType string) {
	for _, link := range engineBinarySymlinks(engineType) {
		if target, err := os.Readlink(link.newname); err == nil && target == link.oldname {
			syscall.Unlink(link.newname)
		}
	}
	for file := range engineScripts(cfg, engineType) {
		os.Remove(file)
	}
}

func engineScripts(cfg *config.CloudConfig, engineType string) map[string]string {
	starter, ok := engineStarters[engineType]
	if !ok || starter.scripts == nil {
		return nil
	}
	return starter.scripts(cfg)
}

// writeEngineConfig writes a configuration file for an engine, creating
// its directory. An empty content removes the file so that a setting
// removed from cloud-config does not linger in the console.
func writeEngineConfig(file string, content []byte) error {
	if len(content) == 0 {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(path.Dir(file), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(file, content, 0600)
}

// engineLogCommand wraps an engine daemon command line so that its output is
// appended to a log file, like user Docker's /var/log/docker.log.
func engineLogCommand(logFile string, command []string) []string {
	quoted := ""
	for _, arg := range command {
		quoted += " " + shellQuote(arg)
	}
	return []string{"bash", "-c", fmt.Sprintf("exec%s >> %s 2>&1", quoted, logFile)}
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
