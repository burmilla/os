package control

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/burmilla/os/cmd/respawn"
	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/log"

	"github.com/codegangsta/cli"
)

// consoleAccessCommands are the commands of the console's respawn.conf which
// give interactive access to the system: sshd and the gettys (with autologin)
// on the ttys. They are named like respawn names its markers (the base name
// of the command).
var consoleAccessCommands = []string{"sshd", "agetty"}

func consoleAccessSubcommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "status",
			Usage:  "print whether interactive access to the console (sshd and tty logins) is enabled or disabled",
			Action: consoleAccessStatus,
		},
		{
			Name:   "disable",
			Usage:  "stop sshd and the tty logins of the console and keep them off across reboots (rancher.console_access: false)",
			Action: consoleAccessDisable,
		},
		{
			Name:   "enable",
			Usage:  "start sshd and the tty logins of the console again (rancher.console_access: true)",
			Action: consoleAccessEnable,
		},
	}
}

func consoleAccessStatus(c *cli.Context) error {
	if config.LoadConfig().Rancher.ConsoleAccess {
		fmt.Println("enabled")
	} else {
		fmt.Println("disabled")
	}
	return nil
}

func consoleAccessDisable(c *cli.Context) error {
	if err := config.Set("rancher.console_access", false); err != nil {
		return err
	}
	if err := disableConsoleAccess(); err != nil {
		return err
	}
	fmt.Println("disabled")
	return nil
}

func consoleAccessEnable(c *cli.Context) error {
	if err := config.Set("rancher.console_access", true); err != nil {
		return err
	}
	// respawn starts the commands again as soon as their markers are gone.
	for _, command := range consoleAccessCommands {
		if err := os.Remove(respawnMarker(command)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	fmt.Println("enabled")
	return nil
}

func respawnMarker(command string) string {
	return filepath.Join(respawn.DisabledDir, command)
}

// writeConsoleAccessMarkers keeps respawn from starting sshd and the gettys.
// The markers live in /run, which the console shares with the host and the
// other system containers, so this works from outside of the console too.
func writeConsoleAccessMarkers() error {
	if err := os.MkdirAll(respawn.DisabledDir, 0755); err != nil {
		return err
	}
	for _, command := range consoleAccessCommands {
		if err := os.WriteFile(respawnMarker(command), nil, 0644); err != nil {
			return err
		}
	}
	return nil
}

// disableConsoleAccess stops sshd and the tty logins of the running console.
// Established SSH sessions are kept, like with ros console sshd disable; tty
// sessions are hung up because an autologin shell would otherwise stay
// available on the ttys.
func disableConsoleAccess() error {
	if err := writeConsoleAccessMarkers(); err != nil {
		return err
	}
	// respawn stops the processes of the disabled commands itself. sshd is
	// also stopped through its pid file, for consoles whose respawn does not
	// handle SIGHUP yet.
	if err := signalRespawn(syscall.SIGHUP); err != nil {
		log.Warnf("Cannot tell respawn of the console to stop the tty logins: %v", err)
	}
	return stopSshd()
}

// signalRespawn sends a signal to the respawn of the console, found through
// the pid file it writes to the shared /run.
func signalRespawn(sig syscall.Signal) error {
	data, err := os.ReadFile(respawn.PidFile)
	if err != nil {
		return err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || pid <= 1 {
		return fmt.Errorf("invalid pid in %s", respawn.PidFile)
	}
	// The pid file survives a console restart, make sure the pid is respawn.
	comm, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return fmt.Errorf("respawn (pid %d) is not running: %v", pid, err)
	}
	if strings.TrimSpace(string(comm)) != "respawn" {
		return fmt.Errorf("pid %d from %s is not respawn", pid, respawn.PidFile)
	}
	return syscall.Kill(pid, sig)
}
