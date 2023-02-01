package sysinit

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/burmilla/os/cmd/control"
	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/compose"
	"github.com/burmilla/os/pkg/docker"
	"github.com/burmilla/os/pkg/log"

	"github.com/docker/libcompose/project/options"
	"golang.org/x/net/context"
)

const (
	systemImagesPreloadDirectory = "/var/lib/rancher/preload/system-docker"
	systemImagesLoadStamp        = "/var/lib/rancher/.sysimages_%s_loaded.done"
)

func hasImage(name string) bool {
	stamp := path.Join(config.StateDir, name)
	if _, err := os.Stat(stamp); os.IsNotExist(err) {
		return false
	}
	return true
}

func getImagesArchive(bootstrap bool) string {
	var archive string
	if bootstrap {
		archive = path.Join(config.ImagesPath, config.InitImages)
	} else {
		archive = path.Join(config.ImagesPath, config.SystemImages)
	}

	return archive
}

func LoadBootstrapImages(cfg *config.CloudConfig) (*config.CloudConfig, error) {
	return loadImages(cfg, true)
}

func LoadSystemImages(cfg *config.CloudConfig) (*config.CloudConfig, error) {
	stamp := fmt.Sprintf(systemImagesLoadStamp, strings.Replace(config.Version, ".", "_", -1))
	if _, err := os.Stat(stamp); os.IsNotExist(err) {
		os.Create(stamp)
		return loadImages(cfg, false)
	}

	log.Infof("Skipped loading system images because %s exists", systemImagesLoadStamp)
	return cfg, nil
}

func loadImages(cfg *config.CloudConfig, bootstrap bool) (*config.CloudConfig, error) {
	archive := getImagesArchive(bootstrap)

	if !hasImage(filepath.Base(archive)) {
		if _, err := os.Stat(archive); os.IsNotExist(err) {
			log.Fatalf("FATAL: Could not load images from %s (file not found)", archive)
		}

		// FixMe: Figure out better logic than just random sleep time
		log.Infof("Waiting 30 seconds before start loading images to containerd")
		time.Sleep(30 * time.Second)

		// client.ImageLoad is an asynchronous operation
		// To ensure the order of execution, use cmd instead of it
		log.Infof("Loading images from %s", archive)
		cmd := exec.Command("/usr/bin/system-docker", "load", "--quiet", "-i", archive)
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("FATAL: Error loading images from %s (%v)\n%s ", archive, err, out)
		}

		log.Infof("Done loading images from %s", archive)
	}

	cmd := exec.Command("/usr/bin/system-docker", "images", "--format", "{{.Repository}}:{{.Tag}}")
	out, _ := cmd.CombinedOutput()
	log.Debugf("Loaded a docker image: %s", out)

	return cfg, nil
}

func SysInit() error {
	cfg := config.LoadConfig()

	if err := control.PreloadImages(docker.NewSystemClient, systemImagesPreloadDirectory); err != nil {
		log.Errorf("Failed to preload System Docker images: %v", err)
	}

	_, err := config.ChainCfgFuncs(cfg,
		config.CfgFuncs{
			{Name: "loadSystemImages", Func: LoadSystemImages},
			{Name: "start project", Func: func(cfg *config.CloudConfig) (*config.CloudConfig, error) {
				p, err := compose.GetProject(cfg, false, true)
				if err != nil {
					return cfg, err
				}
				return cfg, p.Up(context.Background(), options.Up{
					Create: options.Create{
						NoRecreate: true,
					},
					Log: cfg.Rancher.Log,
				})
			}},
			{Name: "sync", Func: func(cfg *config.CloudConfig) (*config.CloudConfig, error) {
				syscall.Sync()
				return cfg, nil
			}},
			{Name: "banner", Func: func(cfg *config.CloudConfig) (*config.CloudConfig, error) {
				log.Infof("BurmillaOS %s started", config.Version)
				return cfg, nil
			}}})
	return err
}

func loadServicesCache() {
	// this code make sure the open-vm-tools, modem-manager... services can be started correct when there is no network
	// make sure the cache directory exist
	if err := os.MkdirAll("/var/lib/rancher/cache/", os.ModeDir|0755); err != nil {
		log.Errorf("Create service cache diretory error: %v", err)
	}

	// move os-services cache file
	if _, err := os.Stat("/usr/share/ros/services-cache"); err == nil {
		files, err := ioutil.ReadDir("/usr/share/ros/services-cache/")
		if err != nil {
			log.Errorf("Read file error: %v", err)
		}
		for _, f := range files {
			err := os.Rename("/usr/share/ros/services-cache/"+f.Name(), "/var/lib/rancher/cache/"+f.Name())
			if err != nil {
				log.Errorf("Rename file error: %v", err)
			}
		}
		if err := os.Remove("/usr/share/ros/services-cache"); err != nil {
			log.Errorf("Remove file error: %v", err)
		}
	}
}

func RunSysInit(c *config.CloudConfig) (*config.CloudConfig, error) {
	loadServicesCache()
	args := append([]string{config.SysInitBin}, os.Args[1:]...)

	cmd := &exec.Cmd{
		Path: config.RosBin,
		Args: args,
	}

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return c, err
	}

	return c, os.Stdin.Close()
}
