package install

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/config/cmdline"
	"github.com/burmilla/os/pkg/log"
	"github.com/burmilla/os/pkg/util"

	"github.com/pkg/errors"
)

const (
	// BootLabel is the label of the traditional separate ext4 boot partition
	BootLabel = "RANCHER_BOOT"
	// EFILabel is the label of the FAT32 EFI system partition created by UEFI
	// installs. FAT labels are limited to 11 characters, so RANCHER_BOOT can
	// not be used there.
	EFILabel = "RANCHER_EFI"
)

// IsEFIFirmware returns true when the running system was booted via UEFI
func IsEFIFirmware() bool {
	_, err := os.Stat("/sys/firmware/efi")
	return err == nil
}

// ResolveLabel returns the device and filesystem type of the filesystem with
// the given label, or empty strings when there is none.
//
// Both the blkid binary and libblkid (in process, through cgo) are used: they
// do not always see the same devices. blkid(8) needs every partition node to
// be present in /dev, which is not guaranteed inside a container - the ros
// upgrade container gets the snapshot of /dev that docker built for it and has
// no udev of its own. Missing the boot partition here is not harmless: the
// caller then falls back to the state partition and writes the new kernel,
// initrd and Syslinux cfgs to a directory which is hidden by the real boot
// partition once the system reboots.
func ResolveLabel(label string) (string, string) {
	d, t, err := util.Blkid(label)
	if err != nil {
		log.Errorf("Failed to run blkid: %s", err)
	}
	if d != "" {
		return d, t
	}

	d = util.ResolveDevice("LABEL=" + label)
	if d == "" {
		return "", ""
	}
	t, err = util.GetFsType(d)
	if err != nil {
		log.Errorf("Failed to get the filesystem type of %s: %s", d, err)
		t = ""
	}
	return d, t
}

// GetBootPartition returns the device and filesystem type of the separate
// boot partition, if there is one - either the traditional ext4 RANCHER_BOOT
// partition or the RANCHER_EFI EFI system partition of a UEFI install.
func GetBootPartition() (string, string) {
	// An installed system records its boot partition on the kernel command
	// line. /proc is the host's one inside the upgrade container, so this is
	// the system's own statement about where its boot files live and it is
	// preferred over guessing from the labels which happen to be around.
	if spec, ok := cmdline.GetCmdline("rancher.state.boot_dev").(string); ok && spec != "" {
		if d := util.ResolveDevice(spec); d != "" {
			t, err := util.GetFsType(d)
			if err != nil {
				log.Errorf("Failed to get the filesystem type of %s: %s", d, err)
				t = ""
			}
			return d, t
		}
		log.Errorf("Could not resolve rancher.state.boot_dev %q", spec)
	}

	for _, label := range []string{BootLabel, EFILabel} {
		if d, t := ResolveLabel(label); d != "" {
			return d, t
		}
	}
	return "", ""
}

func MountDevice(baseName, device, partition string, raw bool) (string, string, error) {
	log.Debugf("mountdevice %s, raw %v", partition, raw)

	fsType := ""
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

		if d, t := GetBootPartition(); d != "" {
			partition = d
			fsType = t
			baseName = filepath.Join(baseName, config.BootDir)
		} else {
			partition = GetStatePartition()
		}
		if partition == "" {
			return device, partition, fmt.Errorf("could not find the partition to mount on %s", baseName)
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
	log.Debugf("mountdevice return2 -> d: %s, p: %s", device, partition)

	// Be explicit about the filesystem type when we know it - mount guessing
	// it wrong (msdos instead of vfat) costs long file name support. Fall back
	// to letting mount detect the type if it does not like ours.
	err := mountPartition(partition, baseName, fsType)
	if err != nil && fsType != "" {
		log.Errorf("%s, retrying without an explicit filesystem type", err)
		err = mountPartition(partition, baseName, "")
	}
	return device, partition, err
}

func mountPartition(partition, target, fsType string) error {
	args := []string{}
	if fsType != "" {
		args = append(args, "-t", fsType)
	}
	args = append(args, partition, target)

	cmd := exec.Command("mount", args...)
	log.Debugf("Run(%v)", cmd)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to mount %s on %s: %s",
			partition, target, strings.TrimSpace(errBuf.String()))
	}
	return nil
}

// VerifyBootMounted checks that the boot directory of the target really is the
// separate boot partition, and not just an empty directory on the state
// partition which the boot partition hides once the system is running.
// Installing there silently looks like a successful upgrade while the
// bootloader keeps loading the previous version.
func VerifyBootMounted(baseName string) error {
	bootDir := filepath.Join(baseName, config.BootDir)
	fsType, err := util.GetMountFsType(bootDir)
	if err != nil {
		return errors.Wrapf(err, "failed to check what is mounted on %s", bootDir)
	}
	if fsType == "" {
		return fmt.Errorf("%s is not a mount point - the boot partition (LABEL=%s) could not be found or mounted, refusing to write the boot files where the firmware will not read them",
			bootDir, EFILabel)
	}
	log.Debugf("VerifyBootMounted: %s is a %s mount", bootDir, fsType)
	return nil
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

// GetPartition returns the device name of the given partition number, taking
// care of the "p" separator needed by devices whose name ends with a digit
func GetPartition(device string, number int) string {
	for _, prefix := range []string{"nvme", "mmcblk", "loop"} {
		if strings.Contains(device, prefix) {
			return fmt.Sprintf("%sp%d", device, number)
		}
	}
	return fmt.Sprintf("%s%d", device, number)
}

func GetDefaultPartition(device string) string {
	return GetPartition(device, 1)
}
