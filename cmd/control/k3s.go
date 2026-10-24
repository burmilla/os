package control

import (
	"fmt"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/log"

	yaml "github.com/cloudfoundry-incubator/candiedyaml"
)

// K3s engine: the os-k3s image ships the k3s binary in /engine and its
// service runs 'ros engine-start k3s'. k3s runs inside the console like user
// Docker does, so kubectl, crictl and ctr (k3s is a multicall binary) work in
// the console and its state lives in /var/lib/rancher/k3s on the state
// partition.
const (
	k3sBinary     = "/usr/bin/k3s"
	k3sConfigFile = "/etc/rancher/k3s/config.yaml"
	k3sLog        = "/var/log/k3s.log"
)

func init() {
	engineStarters[config.EngineK3s] = engineStarter{
		start: startK3s,
		init:  k3sInit,
		links: []symlink{
			{k3sBinary, "/usr/bin/kubectl"},
			{k3sBinary, "/usr/bin/crictl"},
			{k3sBinary, "/usr/bin/ctr"},
		},
	}
}

func startK3s(cfg *config.CloudConfig) error {
	return startInConsole(cfg, config.EngineK3s, cfg.Rancher.K3s.Environment)
}

func k3sInit(cfg *config.CloudConfig, _ []string) error {
	k3sCfg := cfg.Rancher.K3s

	command := k3sCfg.Command
	if command == "" {
		command = "server"
	}
	if command != "server" && command != "agent" {
		return fmt.Errorf("rancher.k3s.command must be server or agent, not %s", command)
	}

	// rancher.k3s.config is written to the k3s configuration file, where
	// each key is a k3s flag without the leading dashes.
	content := []byte{}
	if len(k3sCfg.Config) > 0 {
		var err error
		if content, err = yaml.Marshal(k3sCfg.Config); err != nil {
			return fmt.Errorf("Failed to marshal rancher.k3s.config: %v", err)
		}
	}
	if err := writeEngineConfig(k3sConfigFile, content); err != nil {
		log.Errorf("Failed to write %s: %v", k3sConfigFile, err)
	}

	args := append([]string{k3sBinary, command}, k3sCfg.ExtraArgs...)
	log.Infof("Running %v", args)

	return execCommand(engineLogCommand(k3sLog, args))
}
