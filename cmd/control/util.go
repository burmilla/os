package control

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/log"

	"github.com/pkg/errors"
)

func yes(question string) bool {
	fmt.Printf("%s [y/N]: ", question)
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return strings.ToLower(line[0:1]) == "y"
}

func formatImage(image string, cfg *config.CloudConfig) string {
	domainRegistry := cfg.Rancher.Environment["REGISTRY_DOMAIN"]
	if domainRegistry != "docker.io" && domainRegistry != "" {
		return fmt.Sprintf("%s/%s", domainRegistry, image)
	}
	return image
}

func symLinkEngineBinary() []symlink {
	baseSymlink := []symlink{
		{"/usr/share/ros/os-release", "/usr/lib/os-release"},
		{"/usr/share/ros/os-release", "/etc/os-release"},

		{"/var/lib/rancher/engine/docker", "/usr/bin/docker"},
		{"/var/lib/rancher/engine/dockerd", "/usr/bin/dockerd"},
		{"/var/lib/rancher/engine/docker-init", "/usr/bin/docker-init"},
		{"/var/lib/rancher/engine/docker-proxy", "/usr/bin/docker-proxy"},

		// >= 18.09.0
		{"/var/lib/rancher/engine/containerd", "/usr/bin/containerd"},
		{"/var/lib/rancher/engine/ctr", "/usr/bin/ctr"},
		{"/var/lib/rancher/engine/containerd-shim", "/usr/bin/containerd-shim"},
		{"/var/lib/rancher/engine/runc", "/usr/bin/runc"},

		// < 18.09.0
		{"/var/lib/rancher/engine/docker-containerd", "/usr/bin/docker-containerd"},
		{"/var/lib/rancher/engine/docker-containerd-ctr", "/usr/bin/docker-containerd-ctr"},
		{"/var/lib/rancher/engine/docker-containerd-shim", "/usr/bin/docker-containerd-shim"},
		{"/var/lib/rancher/engine/docker-runc", "/usr/bin/docker-runc"},

		// Docker CLI plugins
		{"/var/lib/rancher/engine-plugins/docker-compose", "/usr/bin/docker-compose"},
		{"/var/lib/rancher/engine-plugins/docker-compose", "/usr/libexec/docker/cli-plugins/docker-compose"},
		{"/var/lib/rancher/engine-plugins/docker-buildx", "/usr/libexec/docker/cli-plugins/docker-buildx"},
	}
	return baseSymlink
}

func checkZfsBackingFS(driver, dir string) error {
	if driver != "zfs" {
		return nil
	}
	for i := 0; i < 4; i++ {
		mountInfo, err := ioutil.ReadFile("/proc/self/mountinfo")
		if err != nil {
			continue
		}
		for _, mount := range strings.Split(string(mountInfo), "\n") {
			if strings.Contains(mount, dir) && strings.Contains(mount, driver) {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
	return errors.Errorf("BackingFS: %s not match storage-driver: %s", dir, driver)
}

func checkGlobalCfg() bool {
	_, err := os.Stat("/proc/1/root/boot/global.cfg")
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}
