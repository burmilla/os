package install

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ReadGlobalCfg(globalCfg string) (string, error) {
	append := ""
	buf, err := ioutil.ReadFile(globalCfg)
	if err != nil {
		return append, err
	}

	s := bufio.NewScanner(bytes.NewReader(buf))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "APPEND") {
			append = strings.TrimSpace(strings.TrimPrefix(line, "APPEND"))
		}
	}
	return append, nil
}

func ReadSyslinuxCfg(currentCfg string) (string, string, error) {
	vmlinuzFile := ""
	initrdFile := ""
	// Need to parse currentCfg for the lines:
	// KERNEL ../vmlinuz-4.9.18-rancher^M
	// INITRD ../initrd-41e02e6-dirty^M
	buf, err := ioutil.ReadFile(currentCfg)
	if err != nil {
		return vmlinuzFile, initrdFile, err
	}

	DIST := filepath.Dir(currentCfg)

	s := bufio.NewScanner(bytes.NewReader(buf))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "KERNEL") {
			vmlinuzFile = strings.TrimSpace(strings.TrimPrefix(line, "KERNEL"))
			vmlinuzFile = filepath.Join(DIST, filepath.Base(vmlinuzFile))
		}
		if strings.HasPrefix(line, "INITRD") {
			initrdFile = strings.TrimSpace(strings.TrimPrefix(line, "INITRD"))
			initrdFile = filepath.Join(DIST, filepath.Base(initrdFile))
		}
	}
	return vmlinuzFile, initrdFile, err
}
