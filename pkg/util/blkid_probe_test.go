package util

import (
	"bytes"
	"encoding/binary"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// blkidValue returns what the system blkid reports for tag on device, or ""
// when blkid is not available.
func blkidValue(t *testing.T, device, tag string) string {
	t.Helper()
	blkid, err := exec.LookPath("blkid")
	if err != nil {
		return ""
	}
	out, err := exec.Command(blkid, "-p", "-o", "value", "-s", tag, device).Output()
	if err != nil {
		t.Logf("blkid -s %s %s: %v", tag, device, err)
		return ""
	}
	return strings.TrimSpace(string(out))
}

func makeImage(t *testing.T, size int) string {
	t.Helper()
	img := filepath.Join(t.TempDir(), "img")
	if err := os.WriteFile(img, make([]byte, size), 0600); err != nil {
		t.Fatal(err)
	}
	return img
}

func probeImage(t *testing.T, img string) *fsProbe {
	t.Helper()
	f, err := os.Open(img)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	return probeFs(f)
}

func TestProbeExt4(t *testing.T) {
	mkfs, err := exec.LookPath("mkfs.ext4")
	if err != nil {
		t.Skip("mkfs.ext4 not available")
	}
	img := makeImage(t, 8<<20)
	if out, err := exec.Command(mkfs, "-q", "-L", "RANCHER_STATE", img).CombinedOutput(); err != nil {
		t.Fatalf("mkfs.ext4: %v: %s", err, out)
	}
	p := probeImage(t, img)
	if p == nil || p.Type != "ext4" || p.Label != "RANCHER_STATE" {
		t.Fatalf("probe = %+v; want ext4 RANCHER_STATE", p)
	}
	if want := blkidValue(t, img, "UUID"); want != "" && p.UUID != want {
		t.Errorf("UUID = %q; blkid says %q", p.UUID, want)
	}
	if got, err := getFsTypeProbe(img); err != nil || got != "ext4" {
		t.Errorf("getFsTypeProbe = %q, %v; want ext4", got, err)
	}
	if _, err := getFsTypeProbe(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("getFsTypeProbe on a missing file should fail")
	}
}

func TestProbeExt2AndExt3(t *testing.T) {
	mkfs, err := exec.LookPath("mke2fs")
	if err != nil {
		t.Skip("mke2fs not available")
	}
	for _, fs := range []string{"ext2", "ext3"} {
		img := makeImage(t, 8<<20)
		if out, err := exec.Command(mkfs, "-q", "-t", fs, "-L", fs, img).CombinedOutput(); err != nil {
			t.Fatalf("mke2fs -t %s: %v: %s", fs, err, out)
		}
		p := probeImage(t, img)
		if p == nil || p.Type != fs || p.Label != fs {
			t.Errorf("probe = %+v; want %s", p, fs)
		}
	}
}

func TestProbeSwap(t *testing.T) {
	mkswap, err := exec.LookPath("mkswap")
	if err != nil {
		t.Skip("mkswap not available")
	}
	img := makeImage(t, 1<<20)
	if out, err := exec.Command(mkswap, "-L", "SWAP", img).CombinedOutput(); err != nil {
		t.Fatalf("mkswap: %v: %s", err, out)
	}
	p := probeImage(t, img)
	if p == nil || p.Type != "swap" || p.Label != "SWAP" {
		t.Fatalf("probe = %+v; want swap SWAP", p)
	}
	if want := blkidValue(t, img, "UUID"); want != "" && p.UUID != want {
		t.Errorf("UUID = %q; blkid says %q", p.UUID, want)
	}
}

// The remaining formats are checked against hand made headers, as the mkfs
// tools are not always available where the tests run.

func TestProbeVfat(t *testing.T) {
	bs := make([]byte, 1<<16)
	bs[0] = 0xEB
	binary.LittleEndian.PutUint16(bs[11:], 512)
	bs[13] = 8
	binary.LittleEndian.PutUint16(bs[14:], 32)
	bs[16] = 2
	binary.LittleEndian.PutUint32(bs[0x43:], 0x1234ABCD)
	copy(bs[0x47:], "BOOT       ")
	copy(bs[0x52:], "FAT32   ")
	bs[510], bs[511] = 0x55, 0xAA
	p := probeFs(bytes.NewReader(bs))
	if p == nil || p.Type != "vfat" || p.Label != "BOOT" || p.UUID != "1234-ABCD" {
		t.Errorf("FAT32 probe = %+v", p)
	}

	bs = make([]byte, 1<<16)
	bs[0] = 0xEB
	binary.LittleEndian.PutUint16(bs[11:], 512)
	bs[13] = 4
	binary.LittleEndian.PutUint16(bs[14:], 1)
	bs[16] = 2
	binary.LittleEndian.PutUint32(bs[0x27:], 0xDEADBEEF)
	copy(bs[0x2b:], "NO NAME    ")
	copy(bs[0x36:], "FAT16   ")
	bs[510], bs[511] = 0x55, 0xAA
	p = probeFs(bytes.NewReader(bs))
	if p == nil || p.Type != "vfat" || p.Label != "" || p.UUID != "DEAD-BEEF" {
		t.Errorf("FAT16 probe = %+v", p)
	}

	// a MBR also ends with 0x55AA but is not a FAT boot sector
	mbr := make([]byte, 1<<16)
	mbr[510], mbr[511] = 0x55, 0xAA
	if p := probeFs(bytes.NewReader(mbr)); p != nil {
		t.Errorf("MBR probed as %+v", p)
	}
}

func TestProbeXfsBtrfsF2fsIso(t *testing.T) {
	uuid := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
	wantUUID := "01020304-0506-0708-090a-0b0c0d0e0f10"

	xfs := make([]byte, 1<<17)
	copy(xfs[0:], "XFSB")
	copy(xfs[32:], uuid)
	copy(xfs[108:], "xfslabel")
	if p := probeFs(bytes.NewReader(xfs)); p == nil || p.Type != "xfs" || p.Label != "xfslabel" || p.UUID != wantUUID {
		t.Errorf("xfs probe = %+v", p)
	}

	btrfs := make([]byte, 1<<17)
	copy(btrfs[0x10000+0x40:], "_BHRfS_M")
	copy(btrfs[0x10000+0x20:], uuid)
	copy(btrfs[0x10000+0x12b:], "btrfslabel")
	if p := probeFs(bytes.NewReader(btrfs)); p == nil || p.Type != "btrfs" || p.Label != "btrfslabel" || p.UUID != wantUUID {
		t.Errorf("btrfs probe = %+v", p)
	}

	f2fs := make([]byte, 1<<17)
	binary.LittleEndian.PutUint32(f2fs[1024:], 0xF2F52010)
	copy(f2fs[1024+0x6c:], uuid)
	for i, c := range "f2fslabel" {
		binary.LittleEndian.PutUint16(f2fs[1024+0x7c+2*i:], uint16(c))
	}
	if p := probeFs(bytes.NewReader(f2fs)); p == nil || p.Type != "f2fs" || p.Label != "f2fslabel" || p.UUID != wantUUID {
		t.Errorf("f2fs probe = %+v", p)
	}

	iso := make([]byte, 1<<17)
	copy(iso[0x8001:], "CD001")
	copy(iso[0x8028:], "BurmillaOS                      ")
	if p := probeFs(bytes.NewReader(iso)); p == nil || p.Type != "iso9660" || p.Label != "BurmillaOS" {
		t.Errorf("iso9660 probe = %+v", p)
	}

	if p := probeFs(bytes.NewReader(make([]byte, 1<<17))); p != nil {
		t.Errorf("zeroes probed as %+v", p)
	}
	if p := probeFs(bytes.NewReader(make([]byte, 100))); p != nil {
		t.Errorf("short device probed as %+v", p)
	}
}

func TestProbePartitions(t *testing.T) {
	disk := make([]byte, 1<<20)
	// protective MBR
	disk[510], disk[511] = 0x55, 0xAA
	// GPT header in LBA 1, entries from LBA 2
	copy(disk[512:], "EFI PART")
	binary.LittleEndian.PutUint64(disk[512+72:], 2)
	binary.LittleEndian.PutUint32(disk[512+80:], 128)
	binary.LittleEndian.PutUint32(disk[512+84:], 128)
	entry := disk[1024:]
	copy(entry[0:], []byte{0xaf, 0x3d, 0xc6, 0x0f, 0x83, 0x84, 0x72, 0x47, 0x8e, 0x79, 0x3d, 0x69, 0xd8, 0x47, 0x7d, 0xe4}) // Linux filesystem
	copy(entry[16:], []byte{0x78, 0x56, 0x34, 0x12, 0xbc, 0x9a, 0xf0, 0xde, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef})
	for i, c := range "RANCHER_STATE" {
		binary.LittleEndian.PutUint16(entry[56+2*i:], uint16(c))
	}
	parts := probePartitions(bytes.NewReader(disk))
	if len(parts) != 1 || parts[0].Number != 1 || parts[0].Label != "RANCHER_STATE" || parts[0].UUID != "12345678-9abc-def0-0123-456789abcdef" {
		t.Errorf("GPT partitions = %+v", parts)
	}

	mbr := make([]byte, 1<<20)
	mbr[510], mbr[511] = 0x55, 0xAA
	binary.LittleEndian.PutUint32(mbr[440:], 0xcafebabe)
	mbr[446+4] = 0x83
	mbr[446+16+4] = 0x82
	parts = probePartitions(bytes.NewReader(mbr))
	if len(parts) != 2 || parts[0].UUID != "cafebabe-01" || parts[1].UUID != "cafebabe-02" || parts[1].Number != 2 {
		t.Errorf("MBR partitions = %+v", parts)
	}

	if got := partitionDevice("sda", 2); got != "sda2" {
		t.Errorf("partitionDevice(sda, 2) = %q", got)
	}
	if got := partitionDevice("nvme0n1", 2); got != "nvme0n1p2" {
		t.Errorf("partitionDevice(nvme0n1, 2) = %q", got)
	}
}

func TestResolveDeviceSpec(t *testing.T) {
	if got := resolveDeviceSpec(""); got != "" {
		t.Errorf("empty spec = %q", got)
	}
	if got := resolveDeviceSpec("/dev/does-not-exist"); got != "/dev/does-not-exist" {
		t.Errorf("plain path = %q", got)
	}
	if got := resolveDeviceSpec("LABEL=surely-no-such-label-here"); got != "" {
		t.Errorf("unknown label = %q", got)
	}
	if got := resolveDeviceSpec("BOGUS=x"); got != "" {
		t.Errorf("unknown tag = %q", got)
	}
	if got := udevEncode("my label/1"); got != "my\\x20label\\x2f1" {
		t.Errorf("udevEncode = %q", got)
	}
}
