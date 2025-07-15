package install

import (
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/burmilla/os/pkg/log"
)

func RunGrub(baseName, device string) error {
	log.Debugf("installGrub baseName: %s , device %s ", baseName, device)

	//grub-install --boot-directory=${baseName}/boot ${device}
	cmd := exec.Command("grub-install", "--target=x86_64-efi", "--efi-directory="+baseName+"/boot/efi", "--boot-directory=/mnt/new_img/boot", device)
	if err := cmd.Run(); err != nil {
		log.Errorf("%s", err)
		return err
	}
	return nil
}

func grubConfig(menu BootVars) error {
	log.Debugf("grubConfig")

	filetmpl, err := template.New("grub2config").Parse(`{{define "grub2menu"}}menuentry "{{.Name}}" {
  set root=(hd0,msdos1)
  linux /{{.bootDir}}vmlinuz-{{.Version}}-rancheros {{.KernelArgs}} {{.Append}}
  initrd /{{.bootDir}}initrd-{{.Version}}-rancheros
}

{{end}}
set default="0"
set timeout="{{.Timeout}}"
{{if .Fallback}}set fallback={{.Fallback}}{{end}}

{{- range .Entries}}
{{template "grub2menu" .}}
{{- end}}

`)
	if err != nil {
		log.Errorf("grub2config %s", err)
		return err
	}

	cfgFile := filepath.Join(menu.BaseName, menu.BootDir+"grub/grub.cfg")
	log.Debugf("grubConfig written to %s", cfgFile)

	f, err := os.Create(cfgFile)
	if err != nil {
		return err
	}
	err = filetmpl.Execute(f, menu)
	if err != nil {
		return err
	}
	return nil
}

func PvGrubConfig(menu BootVars) error {
	log.Debugf("pvGrubConfig")

	filetmpl, err := template.New("grublst").Parse(`{{define "grubmenu"}}
title BurmillaOS {{.Version}}-({{.Name}})
root (hd0)
kernel /${bootDir}vmlinuz-{{.Version}}-rancheros {{.KernelArgs}} {{.Append}}
initrd /${bootDir}initrd-{{.Version}}-rancheros

{{end}}
default 0
timeout {{.Timeout}}
{{if .Fallback}}fallback {{.Fallback}}{{end}}
hiddenmenu

{{- range .Entries}}
{{template "grubmenu" .}}
{{- end}}

`)
	if err != nil {
		log.Errorf("pv grublst: %s", err)

		return err
	}

	cfgFile := filepath.Join(menu.BaseName, menu.BootDir+"grub/menu.lst")
	log.Debugf("grubMenu written to %s", cfgFile)
	f, err := os.Create(cfgFile)
	if err != nil {
		log.Errorf("Create(%s) %s", cfgFile, err)

		return err
	}
	err = filetmpl.Execute(f, menu)
	if err != nil {
		log.Errorf("execute %s", err)
		return err
	}
	return nil
}
