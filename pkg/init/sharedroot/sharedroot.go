package sharedroot

import (
	"os"

	"github.com/burmilla/os/config"
	"github.com/burmilla/os/pkg/init/fsmount"

	"github.com/docker/docker/pkg/mount"
)

func Setup(c *config.CloudConfig) (*config.CloudConfig, error) {
	if c.Rancher.NoSharedRoot {
		return c, nil
	}

	if fsmount.IsInitrd() {
		for _, i := range []string{"/mnt", "/media", "/var/lib/balena-engine"} {
			if err := os.MkdirAll(i, 0755); err != nil {
				return c, err
			}
			if err := mount.Mount("tmpfs", i, "tmpfs", "rw"); err != nil {
				return c, err
			}
			if err := mount.MakeShared(i); err != nil {
				return c, err
			}
		}
		return c, nil
	}

	return c, mount.MakeShared("/")
}
