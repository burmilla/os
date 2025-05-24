package control

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"

	"github.com/burmilla/os/cmd/control/service"
	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/compose"
	"github.com/burmilla/os/pkg/docker"
	"github.com/burmilla/os/pkg/log"
	"github.com/burmilla/os/pkg/util"
	"github.com/burmilla/os/pkg/util/network"

	"github.com/codegangsta/cli"
	"github.com/docker/docker/reference"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/libcompose/project/options"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func engineSubcommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "switch",
			Usage:  "switch user Docker engine without a reboot",
			Action: engineSwitch,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "do not prompt for input",
				},
				cli.BoolFlag{
					Name:  "no-pull",
					Usage: "don't pull console image",
				},
			},
		},
		{
			Name:   "enable",
			Usage:  "set user Docker engine to be switched on next reboot",
			Action: engineEnable,
		},
		{
			Name:  "list",
			Usage: "list available Docker engines (include the Dind engines)",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "update, u",
					Usage: "update engine cache",
				},
			},
			Action: engineList,
		},
	}
}

func engineSwitch(c *cli.Context) error {
	if len(c.Args()) != 1 {
		log.Fatal("Must specify exactly one Docker engine to switch to")
	}
	newEngine := c.Args()[0]

	cfg := config.LoadConfig()

	if newEngine == "latest" {
		engines := availableEngines(cfg, true)
		newEngine = engines[len(engines)-1]
		currentEngine := CurrentEngine()
		if newEngine == currentEngine {
			log.Infof("Latest engine %s is already running", newEngine)
			return nil
		}
		log.Infof("Switching to engine %s", newEngine)
	} else {
		validateEngine(newEngine, cfg)
	}

	project, err := compose.GetProject(cfg, true, false)
	if err != nil {
		log.Fatal(err)
	}

	if err = project.Stop(context.Background(), 10, "engine"); err != nil {
		log.Fatal(err)
	}

	if err = compose.LoadSpecialService(project, cfg, "engine", newEngine); err != nil {
		log.Fatal(err)
	}

	if err = project.Up(context.Background(), options.Up{}, "engine"); err != nil {
		log.Fatal(err)
	}

	if err := config.Set("rancher.engine", newEngine); err != nil {
		log.Errorf("Failed to update rancher.engine: %v", err)
	}

	return nil
}

func engineCreate(c *cli.Context) error {
	name := c.Args()[0]
	sshPort, _ := strconv.Atoi(c.String("ssh-port"))
	if sshPort <= 0 {
		sshPort = randomSSHPort()
	}

	// stage engine service
	cfg := config.LoadConfig()
	var enabledServices []string
	if val, ok := cfg.Rancher.ServicesInclude[name]; !ok || !val {
		cfg.Rancher.ServicesInclude[name] = true
		enabledServices = append(enabledServices, name)
	}

	if len(enabledServices) > 0 {
		if err := compose.StageServices(cfg, enabledServices...); err != nil {
			log.Fatal(err)
		}

		if err := config.Set("rancher.services_include", cfg.Rancher.ServicesInclude); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func engineEnable(c *cli.Context) error {
	if len(c.Args()) != 1 {
		log.Fatal("Must specify exactly one Docker engine to enable")
	}
	newEngine := c.Args()[0]

	cfg := config.LoadConfig()
	validateEngine(newEngine, cfg)

	if err := compose.StageServices(cfg, newEngine); err != nil {
		return err
	}

	if err := config.Set("rancher.docker.engine", newEngine); err != nil {
		log.Errorf("Failed to update 'rancher.docker.engine': %v", err)
	}

	return nil
}

func engineList(c *cli.Context) error {
	cfg := config.LoadConfig()
	engines := availableEngines(cfg, c.Bool("update"))
	currentEngine := CurrentEngine()

	i := 1
	for _, engine := range engines {
		if engine == currentEngine {
			if i == len(engines) {
				fmt.Printf("current  %s (latest)\n", engine)
			} else {
				fmt.Printf("current  %s\n", engine)
			}
		} else {
			if i == len(engines) {
				fmt.Printf("disabled %s (latest)\n", engine)
			} else {
				fmt.Printf("disabled %s\n", engine)
			}
		}
		i++
	}

	// check the dind container
	client, err := docker.NewSystemClient()
	if err != nil {
		log.Warnf("Failed to detect dind: %v", err)
		return nil
	}

	filter := filters.NewArgs()
	filter.Add("label", config.UserDockerLabel)
	opts := types.ContainerListOptions{
		All:    true,
		Filter: filter,
	}
	containers, err := client.ContainerList(context.Background(), opts)
	if err != nil {
		log.Warnf("Failed to detect dind: %v", err)
		return nil
	}
	for _, c := range containers {
		if c.State == "running" {
			fmt.Printf("enabled  %s\n", c.Labels[config.UserDockerLabel])
		} else {
			fmt.Printf("disabled %s\n", c.Labels[config.UserDockerLabel])
		}
	}

	return nil
}

func validateEngine(engine string, cfg *config.CloudConfig) {
	engines := availableEngines(cfg, false)
	if !service.IsLocalOrURL(engine) && !util.Contains(engines, engine) {
		log.Fatalf("%s is not a valid engine", engine)
	}
}

func availableEngines(cfg *config.CloudConfig, update bool) []string {
	if update {
		err := network.UpdateCaches(cfg.Rancher.Repositories.ToArray(), "engines")
		if err != nil {
			log.Debugf("Failed to update engine caches: %v", err)
		}

	}
	engines, err := network.GetEngines(cfg.Rancher.Repositories.ToArray())
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(engines)
	return engines
}

// CurrentEngine gets the name of the docker that's running
func CurrentEngine() (engine string) {
	// sudo system-docker inspect --format "{{.Config.Image}}" docker
	client, err := docker.NewSystemClient()
	if err != nil {
		log.Warnf("Failed to detect current engine: %v", err)
		return
	}
	info, err := client.ContainerInspect(context.Background(), "engine")
	if err != nil {
		log.Warnf("Failed to detect current engine: %v", err)
		return
	}
	// parse image name, then remove os- prefix and the engine suffix
	image, err := reference.ParseNamed(info.Config.Image)
	if err != nil {
		log.Warnf("Failed to detect current engine(%s): %v", info.Config.Image, err)
		return
	}
	if t, ok := image.(reference.NamedTagged); ok {
		tag := t.Tag()
		return tag
	}

	return
}

func preFlightValidate(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("Must specify one engine name")
	}
	name := c.Args()[0]
	if name == "" {
		return errors.New("Must specify one engine name")
	}

	version := c.String("version")
	if version == "" {
		return errors.New("Must specify one engine version")
	}

	authorizedKeys := c.String("authorized-keys")
	if authorizedKeys != "" {
		if _, err := os.Stat(authorizedKeys); os.IsNotExist(err) {
			return errors.New("The authorized-keys should be an exist file, recommended to put in the /opt or /var/lib/rancher directory")
		}
	}

	network := c.String("network")
	if network == "" {
		return errors.New("Must specify network")
	}

	userDefineNetwork, err := CheckUserDefineNetwork(network)
	if err != nil {
		return err
	}

	fixedIP := c.String("fixed-ip")
	if fixedIP == "" {
		return errors.New("Must specify fix ip")
	}

	err = CheckUserDefineIPv4Address(fixedIP, *userDefineNetwork)
	if err != nil {
		return err
	}

	isVersionMatch := false
	for _, v := range config.SupportedDinds {
		if v == version {
			isVersionMatch = true
			break
		}
	}

	if !isVersionMatch {
		return errors.Errorf("Engine version not supported only %v are supported", config.SupportedDinds)
	}

	if c.String("ssh-port") != "" {
		port, err := strconv.Atoi(c.String("ssh-port"))
		if err != nil {
			return errors.Wrap(err, "Failed to convert ssh port to Int")
		}
		if port > 0 {
			addr, err := net.ResolveTCPAddr("tcp", "localhost:"+strconv.Itoa(port))
			if err != nil {
				return errors.Errorf("Failed to resolve tcp addr: %v", err)
			}
			l, err := net.ListenTCP("tcp", addr)
			if err != nil {
				return errors.Errorf("Failed to listen tcp: %v", err)
			}
			defer l.Close()
		}
	}

	return nil
}

func randomSSHPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		log.Errorf("Failed to resolve tcp addr: %v", err)
		return 0
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func CheckUserDefineNetwork(name string) (*types.NetworkResource, error) {
	systemClient, err := docker.NewSystemClient()
	if err != nil {
		return nil, err
	}

	networks, err := systemClient.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	for _, network := range networks {
		if network.Name == name {
			return &network, nil
		}
	}

	return nil, errors.Errorf("Failed to found the user define network: %s", name)
}

func CheckUserDefineIPv4Address(ipv4 string, network types.NetworkResource) error {
	for _, config := range network.IPAM.Config {
		_, ipnet, _ := net.ParseCIDR(config.Subnet)
		if ipnet.Contains(net.ParseIP(ipv4)) {
			return nil
		}
	}
	return errors.Errorf("IP %s is not in the specified cidr", ipv4)
}
