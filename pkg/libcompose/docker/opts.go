package docker

// Parsers for the restart policy and device strings of a compose file. They
// were taken from github.com/docker/docker/runconfig/opts, which also pulls
// in docker's signal table and with it syscall.SIGUNUSED, which does not
// exist on loong64.

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/engine-api/types/container"
)

// parseRestartPolicy returns the parsed policy or an error indicating what is incorrect
func parseRestartPolicy(policy string) (container.RestartPolicy, error) {
	p := container.RestartPolicy{}

	if policy == "" {
		return p, nil
	}

	var (
		parts = strings.Split(policy, ":")
		name  = parts[0]
	)

	p.Name = name
	switch name {
	case "always", "unless-stopped":
		if len(parts) > 1 {
			return p, fmt.Errorf("maximum restart count not valid with restart policy of \"%s\"", name)
		}
	case "no":
		// do nothing
	case "on-failure":
		if len(parts) > 2 {
			return p, fmt.Errorf("restart count format is not valid, usage: 'on-failure:N' or 'on-failure'")
		}
		if len(parts) == 2 {
			count, err := strconv.Atoi(parts[1])
			if err != nil {
				return p, err
			}

			p.MaximumRetryCount = count
		}
	default:
		return p, fmt.Errorf("invalid restart policy %s", name)
	}

	return p, nil
}

// parseDevice parses a device mapping string to a container.DeviceMapping struct
func parseDevice(device string) (container.DeviceMapping, error) {
	src := ""
	dst := ""
	permissions := "rwm"
	arr := strings.Split(device, ":")
	switch len(arr) {
	case 3:
		permissions = arr[2]
		fallthrough
	case 2:
		if validDeviceMode(arr[1]) {
			permissions = arr[1]
		} else {
			dst = arr[1]
		}
		fallthrough
	case 1:
		src = arr[0]
	default:
		return container.DeviceMapping{}, fmt.Errorf("invalid device specification: %s", device)
	}

	if dst == "" {
		dst = src
	}

	return container.DeviceMapping{
		PathOnHost:        src,
		PathInContainer:   dst,
		CgroupPermissions: permissions,
	}, nil
}

// validDeviceMode checks if the mode for device is valid or not.
// Valid mode is a composition of r (read), w (write), and m (mknod).
func validDeviceMode(mode string) bool {
	var legalDeviceMode = map[rune]bool{
		'r': true,
		'w': true,
		'm': true,
	}
	if mode == "" {
		return false
	}
	for _, c := range mode {
		if !legalDeviceMode[c] {
			return false
		}
		legalDeviceMode[c] = false
	}
	return true
}
