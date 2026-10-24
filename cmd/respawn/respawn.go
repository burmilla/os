package respawn

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/log"

	"github.com/codegangsta/cli"
)

// DisabledDir holds markers named after a command (e.g. "sshd"). While a
// marker exists the command is not (re)started. It lives in /run so every
// command is enabled again after a reboot.
const DisabledDir = "/run/respawn.disabled"

// PidFile is where respawn writes its PID, so that other containers which
// share /run (e.g. ros console access disable) can send it SIGHUP after
// adding a marker to DisabledDir. respawn then stops the processes of the
// disabled commands.
const PidFile = "/run/respawn.pid"

type process struct {
	*os.Process
	command string
}

var (
	running     = true
	processes   = map[int]process{}
	processLock = sync.Mutex{}
)

func Main() {
	log.InitLogger()
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
	app := cli.NewApp()

	app.Name = os.Args[0]
	app.Usage = fmt.Sprintf("%s BurmillaOS\nbuilt: %s", app.Name, config.BuildDate)
	app.Version = config.Version
	app.Author = "Project Burmilla\n\tRancher Labs, Inc."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Optional config file to load",
		},
	}
	app.Action = run

	log.Infof("%s, %s", app.Usage, app.Version)
	fmt.Printf("%s, %s", app.Usage, app.Version)

	app.Run(os.Args)
}

func setupSignals() {
	sigtermChan := make(chan os.Signal, 1)
	signal.Notify(sigtermChan, syscall.SIGTERM)
	go func() {
		for range sigtermChan {
			termPids()
		}
	}()
	sighupChan := make(chan os.Signal, 1)
	signal.Notify(sighupChan, syscall.SIGHUP)
	go func() {
		for range sighupChan {
			termDisabled()
		}
	}()
}

func run(c *cli.Context) error {
	setupSignals()
	if err := ioutil.WriteFile(PidFile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644); err != nil {
		log.Errorf("Failed to write %s: %v", PidFile, err)
	}

	var stream io.Reader = os.Stdin
	var err error

	inputFileName := c.String("file")

	if inputFileName != "" {
		stream, err = os.Open(inputFileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	input, err := ioutil.ReadAll(stream)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")
	doneChannel := make(chan string, len(lines))

	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.Index(strings.TrimSpace(line), "#") == 0 {
			continue
		}
		go execute(line, doneChannel)
	}

	for i := 0; i < len(lines); i++ {
		line := <-doneChannel
		log.Infof("FINISHED: %s", line)
		fmt.Printf("FINISHED: %s", line)
	}
	return nil
}

func addProcess(p *os.Process, command string) {
	processLock.Lock()
	defer processLock.Unlock()
	processes[p.Pid] = process{Process: p, command: command}
}

func removeProcess(process *os.Process) {
	processLock.Lock()
	defer processLock.Unlock()
	delete(processes, process.Pid)
}

func termPids() {
	running = false
	processLock.Lock()
	defer processLock.Unlock()

	for _, process := range processes {
		log.Infof("sending SIGTERM to %d", process.Pid)
		process.Signal(syscall.SIGTERM)
	}
}

// termDisabled stops the running processes of the commands which have been
// disabled since they were started (SIGHUP). The command is not started
// again until its marker is removed. Processes which do not exit on SIGTERM
// within a few seconds (e.g. an interactive shell started by agetty) are
// killed.
func termDisabled() {
	processLock.Lock()
	defer processLock.Unlock()

	for _, p := range processes {
		if !Disabled(p.command) {
			continue
		}
		log.Infof("%s : disabled, sending SIGTERM to %d", p.command, p.Pid)
		if err := p.Signal(syscall.SIGTERM); err != nil {
			continue
		}
		go killAfter(p, 5*time.Second)
	}
}

func killAfter(p process, d time.Duration) {
	time.Sleep(d)
	processLock.Lock()
	defer processLock.Unlock()
	if current, ok := processes[p.Pid]; !ok || current.Process != p.Process {
		return
	}
	log.Infof("%s : %d did not exit, sending SIGKILL", p.command, p.Pid)
	p.Signal(syscall.SIGKILL)
}

func execute(line string, doneChannel chan string) {
	defer func() { doneChannel <- line }()

	start := time.Now()
	count := 0

	args := strings.Split(line, " ")

	for {
		if waitWhileDisabled(args[0]) {
			// Do not count the pause as a fast restart.
			count = 0
			start = time.Now()
		}
		if !running {
			break
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}

		if err := cmd.Start(); err == nil {
			addProcess(cmd.Process, args[0])
			if err = cmd.Wait(); err != nil {
				log.Errorf("Wait cmd to exit: %s, err: %v", line, err)
			}
			removeProcess(cmd.Process)
		} else {
			log.Errorf("Start cmd: %s, err: %v", line, err)
		}

		if !running {
			log.Infof("%s : not restarting, exiting", line)
			break
		}

		count++

		if count > 10 {
			if time.Now().Sub(start) <= (1 * time.Second) {
				log.Errorf("%s : restarted too fast, not executing", line)
				break
			}

			count = 0
			start = time.Now()
		}
	}
}

// Disabled reports whether starting a command is disabled.
func Disabled(command string) bool {
	_, err := os.Stat(filepath.Join(DisabledDir, filepath.Base(command)))
	return err == nil
}

// waitWhileDisabled blocks while the command is disabled and reports
// whether it had to wait.
func waitWhileDisabled(command string) bool {
	waited := false
	for running && Disabled(command) {
		if !waited {
			log.Infof("%s : disabled, waiting for %s to be removed", command, filepath.Join(DisabledDir, filepath.Base(command)))
			waited = true
		}
		time.Sleep(time.Second)
	}
	return waited
}
