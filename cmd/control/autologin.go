package control

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/log"

	"github.com/codegangsta/cli"
)

func AutologinMain() {
	log.InitLogger()
	app := cli.NewApp()

	app.Name = os.Args[0]
	app.Usage = "autologin console"
	app.Version = config.Version
	app.Author = "Project Burmilla\n\tRancher Labs, Inc."
	app.Email = "burmilla@localhost.local"
	app.EnableBashCompletion = true
	app.Action = autologinAction
	app.HideHelp = true
	app.Run(os.Args)
}

func autologinAction(c *cli.Context) error {
	cmd := exec.Command("/bin/stty", "sane")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	usertty := ""
	user := "root"
	if c.NArg() > 0 {
		usertty = c.Args().Get(0)
		s := strings.SplitN(usertty, ":", 2)
		user = s[0]
	}

	mode := filepath.Base(os.Args[0])
	console := CurrentConsole()

	cfg := config.LoadConfig()

	loginBin := ""
	args := []string{}
	if console == "centos" || console == "fedora" ||
		mode == "recovery" {
		// For some reason, centos and fedora ttyS0 and tty1 don't work with `login -f rancher`
		// until I make time to read their source, lets just give us a way to get work done
		loginBin = "bash"
		args = append(args, "--login")
		if mode == "recovery" {
			os.Setenv("PROMPT_COMMAND", `echo "[`+fmt.Sprintf("Recovery console %s@%s:${PWD}", user, cfg.Hostname)+`]"`)
		}
	} else {
		loginBin = "login"
		args = append(args, "-f", user)
		// TODO: add a PROMPT_COMMAND if we haven't switch-rooted
	}

	loginBinPath, err := exec.LookPath(loginBin)
	if err != nil {
		fmt.Printf("error finding %s in path: %s", cmd.Args[0], err)
		return err
	}
	os.Setenv("TERM", "linux")

	// login from util-linux (Debian 13 and later consoles) must run as the
	// session leader owning the tty: when run as a child process instead,
	// its setpgrp() puts it into a background process group and the kernel
	// then stops it with SIGTTOU on the first tcsetattr(), leaving the
	// console hung. Replace this process (exec'd by agetty as session
	// leader) with the login program, as getty implementations normally do.
	if err := syscall.Exec(loginBinPath, append([]string{loginBin}, args...), os.Environ()); err != nil {
		log.Errorf("\nError starting %s: %s", loginBin, err)
		return err
	}
	return nil
}
