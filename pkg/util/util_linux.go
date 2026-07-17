//go:build linux
// +build linux

package util

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/SvenDowideit/cpuid"
	"github.com/docker/docker/pkg/mount"
)

func mountProc() error {
	if _, err := os.Stat("/proc/self/mountinfo"); os.IsNotExist(err) {
		if _, err := os.Stat("/proc"); os.IsNotExist(err) {
			if err = os.Mkdir("/proc", 0755); err != nil {
				return err
			}
		}

		if err := syscall.Mount("none", "/proc", "proc", 0, ""); err != nil {
			return err
		}
	}

	return nil
}

func Mount(device, target, fsType, options string) error {
	if err := mountProc(); err != nil {
		return nil
	}

	bindMount := false
	for _, v := range strings.Split(options, ",") {
		if v == "bind" {
			bindMount = true
			break
		}
	}

	if bindMount {
		deviceInfo, err := os.Stat(device)
		if err != nil {
			return err
		}
		mode := deviceInfo.Mode()

		switch {
		case mode.IsDir():
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case mode.IsRegular():
			err := os.MkdirAll(filepath.Dir(target), 0755)
			if err != nil {
				return err
			}
			file, err := os.OpenFile(target, os.O_CREATE, mode&os.ModePerm)
			if err != nil {
				return err
			}
			if err := file.Close(); err != nil {
				return err
			}
		default:
			return os.ErrInvalid
		}
	} else {
		err := os.MkdirAll(target, 0755)
		if err != nil {
			return err
		}

		if fsType == "auto" || fsType == "" {
			inferredType, err := GetFsType(device)
			if err != nil {
				return err
			}
			fsType = inferredType
		}
	}

	return mount.Mount(device, target, fsType, options)
}

func Unmount(target string) error {
	return mount.Unmount(target)
}

// GetMountFsType returns the filesystem type of the topmost filesystem mounted
// on target or an empty string when nothing is mounted there.
func GetMountFsType(target string) (string, error) {
	mounts, err := mount.GetMounts()
	if err != nil {
		return "", err
	}

	fsType := ""
	for _, m := range mounts {
		// Mounts are listed in mount order, so the last entry matching
		// target is the filesystem which is currently visible on it.
		if m.Mountpoint == target {
			fsType = m.Fstype
		}
	}

	return fsType, nil
}

// blkidTag returns the value of the KEY="value" tag from a single line of
// blkid output. Only whole tag names match, so looking for TYPE does not
// return the value of SEC_TYPE - blkid prints `SEC_TYPE="msdos"` before
// `TYPE="vfat"` for FAT filesystems, and reporting a FAT32 EFI system
// partition as "msdos" makes it get mounted without long file name support.
func blkidTag(line, key string) string {
	tag := key + `="`
	for offset := 0; ; {
		i := strings.Index(line[offset:], tag)
		if i < 0 {
			return ""
		}
		i += offset
		offset = i + len(tag)
		// the tag name has to start the line or follow a separator, otherwise
		// we matched the tail of a longer tag name
		if i != 0 && line[i-1] != ' ' {
			continue
		}
		if end := strings.Index(line[offset:], `"`); end >= 0 {
			return line[offset : offset+end]
		}
		return ""
	}
}

func Blkid(label string) (deviceName, deviceType string, err error) {
	// Not all blkid's have `blkid -L label (see busybox/alpine)
	cmd := exec.Command("blkid")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return
	}
	r := bytes.NewReader(out)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if blkidTag(line, "LABEL") != label {
			continue
		}
		d := strings.Split(line, ":")
		deviceName = d[0]
		deviceType = blkidTag(line, "TYPE")
		return
	}
	return
}

func BlkidType(deviceType string) (deviceNames []string, err error) {
	// Not all blkid's have `blkid -L label (see busybox/alpine)
	cmd := exec.Command("blkid")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(out)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if blkidTag(line, "TYPE") != deviceType {
			continue
		}
		d := strings.Split(line, ":")
		deviceName := d[0]
		deviceNames = append(deviceNames, deviceName)
	}
	return deviceNames, nil
}

// GetHypervisor tries to detect if we're running in a VM, and returns a string for its type
func GetHypervisor() string {
	return cpuid.CPU.HypervisorName
}
