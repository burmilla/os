//go:build linux
// +build linux

package util

import "testing"

func TestBlkidTag(t *testing.T) {
	// a FAT32 EFI system partition as blkid prints it - SEC_TYPE comes before
	// TYPE, and LABEL_FATBOOT before LABEL
	esp := `/dev/sda1: SEC_TYPE="msdos" LABEL_FATBOOT="RANCHER_EFI" LABEL="RANCHER_EFI" UUID="7152-C834" TYPE="vfat" PARTLABEL="ESP" PARTUUID="24bf5701"`
	state := `/dev/sda2: LABEL="RANCHER_STATE" UUID="a0c3c1e9-4644-49c3-b054-a51d1743bec6" TYPE="ext4" PARTUUID="411888c3"`

	tests := []struct {
		line, key, expected string
	}{
		{esp, "TYPE", "vfat"},
		{esp, "SEC_TYPE", "msdos"},
		{esp, "LABEL", "RANCHER_EFI"},
		{esp, "LABEL_FATBOOT", "RANCHER_EFI"},
		{esp, "PARTLABEL", "ESP"},
		{esp, "UUID", "7152-C834"},
		{esp, "MISSING", ""},
		{state, "TYPE", "ext4"},
		{state, "LABEL", "RANCHER_STATE"},
		{state, "SEC_TYPE", ""},
		{state, "LABEL_FATBOOT", ""},
	}

	for _, tc := range tests {
		if actual := blkidTag(tc.line, tc.key); actual != tc.expected {
			t.Errorf("blkidTag(%q) == %q, expected %q", tc.key, actual, tc.expected)
		}
	}
}
