package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/log"
	"github.com/burmilla/os/pkg/util"
)

// IsEFIFirmware returns true when the running system was booted via UEFI
func IsEFIFirmware() bool {
	_, err := os.Stat("/sys/firmware/efi")
	return err == nil
}

// GetBootPartition returns the device and filesystem type of the separate
// boot partition, if there is one. RANCHER_BOOT is the traditional ext4
// boot partition label, RANCHER_EFI is the FAT32 EFI system partition
// created by UEFI installs (FAT labels are limited to 11 characters, so
// RANCHER_BOOT can not be used there).
func GetBootPartition() (string, string) {
	for _, label := range []string{"RANCHER_BOOT", "RANCHER_EFI"} {
		d, t, err := util.Blkid(label)
		if err != nil {
			log.Errorf("Failed to run blkid: %s", err)
			continue
		}
		if d != "" {
			return d, t
		}
	}
	return "", ""
}

func MountDevice(baseName, device, partition string, raw bool) (string, string, error) {
	log.Debugf("mountdevice %s, raw %v", partition, raw)

	if partition == "" {
		if raw {
			log.Debugf("util.Mount (raw) %s, %s", partition, baseName)

			cmd := exec.Command("lsblk", "-no", "pkname", partition)
			log.Debugf("Run(%v)", cmd)
			cmd.Stderr = os.Stderr
			device := ""
			// TODO: out can == "" - this is used to "detect software RAID" which is terrible
			if out, err := cmd.Output(); err == nil {
				device = "/dev/" + strings.TrimSpace(string(out))
			}

			log.Debugf("mountdevice return -> d: %s, p: %s", device, partition)
			return device, partition, util.Mount(partition, baseName, "", "")
		}

		//rootfs := partition
		// Don't use ResolveDevice - it can fail, whereas `blkid -L LABEL` works more often

		if d, _ := GetBootPartition(); d != "" {
			partition = d
			baseName = filepath.Join(baseName, config.BootDir)
		} else {
			partition = GetStatePartition()
		}
		cmd := exec.Command("lsblk", "-no", "pkname", partition)
		log.Debugf("Run(%v)", cmd)
		cmd.Stderr = os.Stderr
		// TODO: out can == "" - this is used to "detect software RAID" which is terrible
		if out, err := cmd.Output(); err == nil {
			device = "/dev/" + strings.TrimSpace(string(out))
		}
	}
	os.MkdirAll(baseName, 0755)
	cmd := exec.Command("mount", partition, baseName)
	//cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	log.Debugf("mountdevice return2 -> d: %s, p: %s", device, partition)
	return device, partition, cmd.Run()
}

func GetStatePartition() string {
	cfg := config.LoadConfig()

	if dev := util.ResolveDevice(cfg.Rancher.State.Dev); dev != "" {
		// try the rancher.state.dev setting
		return dev
	}
	d, _, err := util.Blkid("RANCHER_STATE")
	if err != nil {
		log.Errorf("Failed to run blkid: %s", err)
	}
	return d
}

// GetPartition returns the device name of the given partition number,
// taking care of the "p" separator needed by nvme devices
func GetPartition(device string, number int) string {
	if strings.Contains(device, "nvme") {
		return fmt.Sprintf("%sp%d", device, number)
	}
	return fmt.Sprintf("%s%d", device, number)
}

func GetDefaultPartition(device string) string {
	return GetPartition(device, 1)
}
