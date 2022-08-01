package install

import (
	"os/exec"

	"github.com/burmilla/os/pkg/log"
)

func RunRefind(device string) error {
	log.Debugf("installRefind")

	cmd := exec.Command("refind-install", "--usedefault", device)
	if err := cmd.Run(); err != nil {
		log.Errorf("%s", err)
		return err
	}
	return nil
}
