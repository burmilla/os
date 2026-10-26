package control

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/log"
)

// QEMU engine: the os-qemu image contains QEMU and libvirt, and its service
// runs 'ros engine-start qemu'. Unlike Docker and K3s, which are static
// binaries that run inside the console, libvirt and QEMU need their own root
// filesystem, so the daemons run in the engine container (sharing the host
// network and PID namespaces) and the console gets wrapper scripts that run
// virsh and friends there with system-docker exec. The libvirt socket is in
// /var/run, which the engine container shares with the host.
const (
	qemuDefaultCommand = "libvirtd"
	virtlogd           = "virtlogd"
)

var qemuDefaultTools = []string{"virsh", "virt-install", "qemu-img"}

func init() {
	engineStarters[config.EngineQemu] = engineStarter{
		start:   startQemu,
		init:    qemuInit,
		scripts: qemuScripts,
	}
}

func startQemu(cfg *config.CloudConfig) error {
	qemuCfg := cfg.Rancher.Qemu

	for _, env := range qemuCfg.Environment {
		if err := putEnv(env); err != nil {
			log.Error(err)
		}
	}

	// libvirtd needs virtlogd for the consoles of the virtual machines; it is
	// started by systemd socket activation elsewhere, here start it ourselves.
	if _, err := exec.LookPath(virtlogd); err == nil {
		cmd := exec.Command(virtlogd, "--daemon")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Errorf("Failed to start %s: %v", virtlogd, err)
		}
	}

	command := qemuCfg.Command
	if command == "" {
		command = qemuDefaultCommand
	}
	args := append([]string{command}, qemuCfg.ExtraArgs...)
	log.Infof("Running %v", args)

	return execCommand(args)
}

// qemuInit would run inside the console, but the QEMU engine runs in its own
// container.
func qemuInit(_ *config.CloudConfig, _ []string) error {
	return fmt.Errorf("the %s engine runs in its own container, use 'ros engine-start %s'", config.EngineQemu, config.EngineQemu)
}

// qemuScripts returns the console wrapper scripts for the engine's tools.
func qemuScripts(cfg *config.CloudConfig) map[string]string {
	tools := cfg.Rancher.Qemu.Tools
	if len(tools) == 0 {
		tools = qemuDefaultTools
	}

	scripts := map[string]string{}
	for _, tool := range tools {
		tool = path.Base(tool)
		if tool == "" || tool == "." || tool == "/" {
			continue
		}
		scripts[path.Join("/usr/bin", tool)] = fmt.Sprintf(`#!/bin/sh
# BurmillaOS: %s runs in the %s engine container
if [ -t 0 ]; then tty=-t; fi
exec /usr/bin/system-docker exec -i $tty %s %s "$@"
`, tool, config.EngineQemu, config.EngineQemu, tool)
	}
	return scripts
}

func putEnv(env string) error {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return os.Setenv(env[:i], env[i+1:])
		}
	}
	return fmt.Errorf("invalid environment variable %q, expected NAME=value", env)
}
