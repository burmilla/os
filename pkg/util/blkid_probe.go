package util

// A small pure Go replacement for the parts of libblkid that ros uses
// (blkid_evaluate_spec and blkid_get_tag_value TYPE). It is used when ros is
// built without cgo, which is the case for cross builds (arm64, riscv64 and
// loong64 are built on an amd64 host without a cross C tool chain). Native
// builds keep using libblkid (cutil.go).
//
// It recognises the file systems BurmillaOS installs to or boots from: ext2,
// ext3, ext4, xfs, btrfs, f2fs, vfat, iso9660, swap, and NTFS and exFAT by
// type only. LABEL=, UUID=, PARTLABEL= and PARTUUID= specs are resolved with
// the udev symlinks under /dev/disk when they exist and otherwise by probing
// the block devices listed in /proc/partitions.

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf16"
)

// fsProbe is what a file system probe found on a device.
type fsProbe struct {
	Type  string
	Label string
	UUID  string
}

// readAt reads exactly len(buf) bytes from r at off. A short read is an
// error: the device is smaller than the structure we are looking for.
func readAt(r io.ReaderAt, off int64, size int) ([]byte, error) {
	buf := make([]byte, size)
	n, err := r.ReadAt(buf, off)
	if n == size {
		return buf, nil
	}
	if err == nil {
		err = io.ErrUnexpectedEOF
	}
	return nil, err
}

func cString(b []byte) string {
	if i := bytes.IndexByte(b, 0); i >= 0 {
		b = b[:i]
	}
	return string(b)
}

func utf16String(b []byte) string {
	u := make([]uint16, 0, len(b)/2)
	for i := 0; i+1 < len(b); i += 2 {
		c := binary.LittleEndian.Uint16(b[i:])
		if c == 0 {
			break
		}
		u = append(u, c)
	}
	return string(utf16.Decode(u))
}

// uuidString formats a 16 byte UUID the way libblkid does for ext*, xfs,
// btrfs, f2fs and swap: as a plain big endian byte string.
func uuidString(b []byte) string {
	if bytes.Equal(b, make([]byte, 16)) {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// guidString formats a GPT GUID, whose first three fields are little endian.
func guidString(b []byte) string {
	if bytes.Equal(b, make([]byte, 16)) {
		return ""
	}
	return fmt.Sprintf("%08x-%04x-%04x-%x-%x",
		binary.LittleEndian.Uint32(b[0:4]),
		binary.LittleEndian.Uint16(b[4:6]),
		binary.LittleEndian.Uint16(b[6:8]),
		b[8:10], b[10:16])
}

// probeFs returns the file system found on r, or nil when none is recognised.
func probeFs(r io.ReaderAt) *fsProbe {
	probes := []func(io.ReaderAt) *fsProbe{
		probeExt,
		probeXfs,
		probeBtrfs,
		probeF2fs,
		probeVfat,
		probeNtfsExfat,
		probeIso9660,
		probeSwap,
	}
	for _, probe := range probes {
		if p := probe(r); p != nil {
			return p
		}
	}
	return nil
}

func probeExt(r io.ReaderAt) *fsProbe {
	const (
		extSuperOff   = 1024
		extSuperSize  = 1024
		extMagic      = 0xEF53
		compatJournal = 0x0004
		// features that only ext4 has, from blkid's superblocks/ext.c
		roCompatExt4 = 0x0008 | 0x0010 | 0x0020 | 0x0040 | 0x0400
		incompatExt4 = 0x0040 | 0x0080 | 0x0100 | 0x0200
	)
	sb, err := readAt(r, extSuperOff, extSuperSize)
	if err != nil || binary.LittleEndian.Uint16(sb[0x38:]) != extMagic {
		return nil
	}
	compat := binary.LittleEndian.Uint32(sb[0x5c:])
	incompat := binary.LittleEndian.Uint32(sb[0x60:])
	roCompat := binary.LittleEndian.Uint32(sb[0x64:])
	p := &fsProbe{
		UUID:  uuidString(sb[0x68:0x78]),
		Label: cString(sb[0x78:0x88]),
	}
	switch {
	case incompat&incompatExt4 != 0 || roCompat&roCompatExt4 != 0:
		p.Type = "ext4"
	case compat&compatJournal != 0:
		p.Type = "ext3"
	default:
		p.Type = "ext2"
	}
	return p
}

func probeXfs(r io.ReaderAt) *fsProbe {
	sb, err := readAt(r, 0, 128)
	if err != nil || string(sb[0:4]) != "XFSB" {
		return nil
	}
	return &fsProbe{
		Type:  "xfs",
		UUID:  uuidString(sb[32:48]),
		Label: cString(sb[108:120]),
	}
}

func probeBtrfs(r io.ReaderAt) *fsProbe {
	const superOff = 0x10000
	sb, err := readAt(r, superOff, 0x12b+256)
	if err != nil || string(sb[0x40:0x48]) != "_BHRfS_M" {
		return nil
	}
	return &fsProbe{
		Type:  "btrfs",
		UUID:  uuidString(sb[0x20:0x30]),
		Label: cString(sb[0x12b : 0x12b+256]),
	}
}

func probeF2fs(r io.ReaderAt) *fsProbe {
	const superOff = 1024
	sb, err := readAt(r, superOff, 0x7c+512)
	if err != nil || binary.LittleEndian.Uint32(sb[0:]) != 0xF2F52010 {
		return nil
	}
	return &fsProbe{
		Type:  "f2fs",
		UUID:  uuidString(sb[0x6c:0x7c]),
		Label: utf16String(sb[0x7c : 0x7c+512]),
	}
}

func probeVfat(r io.ReaderAt) *fsProbe {
	bs, err := readAt(r, 0, 512)
	if err != nil || bs[510] != 0x55 || bs[511] != 0xAA {
		return nil
	}
	// a FAT boot sector starts with a jump instruction and has a sector size
	// and a cluster size that are powers of two
	if bs[0] != 0xEB && bs[0] != 0xE9 {
		return nil
	}
	sectorSize := binary.LittleEndian.Uint16(bs[11:])
	clusterSize := bs[13]
	if sectorSize == 0 || sectorSize&(sectorSize-1) != 0 || clusterSize == 0 || clusterSize&(clusterSize-1) != 0 {
		return nil
	}
	var serial uint32
	var label string
	switch {
	case string(bs[0x52:0x5a]) == "FAT32   ":
		serial = binary.LittleEndian.Uint32(bs[0x43:])
		label = string(bs[0x47:0x52])
	case string(bs[0x36:0x3e]) == "FAT12   " || string(bs[0x36:0x3e]) == "FAT16   " || string(bs[0x36:0x3e]) == "FAT     ":
		serial = binary.LittleEndian.Uint32(bs[0x27:])
		label = string(bs[0x2b:0x36])
	default:
		// no file system type string: only accept it when the FAT geometry
		// makes sense (reserved sectors and number of FATs set)
		if binary.LittleEndian.Uint16(bs[14:]) == 0 || bs[16] == 0 {
			return nil
		}
		serial = binary.LittleEndian.Uint32(bs[0x27:])
		label = string(bs[0x2b:0x36])
	}
	label = strings.TrimRight(label, " ")
	if label == "NO NAME" {
		label = ""
	}
	return &fsProbe{
		Type:  "vfat",
		UUID:  fmt.Sprintf("%04X-%04X", serial>>16, serial&0xffff),
		Label: label,
	}
}

func probeNtfsExfat(r io.ReaderAt) *fsProbe {
	bs, err := readAt(r, 0, 512)
	if err != nil {
		return nil
	}
	switch string(bs[3:11]) {
	case "NTFS    ":
		return &fsProbe{Type: "ntfs"}
	case "EXFAT   ":
		return &fsProbe{Type: "exfat"}
	}
	return nil
}

func probeIso9660(r io.ReaderAt) *fsProbe {
	const pvdOff = 0x8000
	pvd, err := readAt(r, pvdOff, 0x28+32)
	if err != nil || string(pvd[1:6]) != "CD001" {
		return nil
	}
	return &fsProbe{
		Type:  "iso9660",
		Label: strings.TrimRight(string(pvd[0x28:0x28+32]), " "),
	}
}

func probeSwap(r io.ReaderAt) *fsProbe {
	// the signature sits at the end of the first page; the page size
	// depends on the architecture the swap area was made on
	for _, pageSize := range []int64{4096, 8192, 16384, 65536} {
		sig, err := readAt(r, pageSize-10, 10)
		if err != nil {
			return nil
		}
		if string(sig) != "SWAPSPACE2" && string(sig) != "SWAP-SPACE" {
			continue
		}
		hdr, err := readAt(r, 1024, 1024)
		if err != nil {
			return nil
		}
		return &fsProbe{
			Type:  "swap",
			UUID:  uuidString(hdr[12:28]),
			Label: cString(hdr[28:44]),
		}
	}
	return nil
}

// partition is a partition table entry of a disk.
type partition struct {
	Number int
	UUID   string // PARTUUID
	Label  string // PARTLABEL
}

// probePartitions lists the partitions of a disk from its GPT or, failing
// that, its MBR.
func probePartitions(r io.ReaderAt) []partition {
	for _, sectorSize := range []int64{512, 4096} {
		if parts := probeGpt(r, sectorSize); parts != nil {
			return parts
		}
	}
	return probeMbr(r)
}

func probeGpt(r io.ReaderAt, sectorSize int64) []partition {
	hdr, err := readAt(r, sectorSize, 92)
	if err != nil || string(hdr[0:8]) != "EFI PART" {
		return nil
	}
	entriesLBA := binary.LittleEndian.Uint64(hdr[72:])
	numEntries := binary.LittleEndian.Uint32(hdr[80:])
	entrySize := binary.LittleEndian.Uint32(hdr[84:])
	if entrySize < 128 || numEntries == 0 || numEntries > 1024 {
		return nil
	}
	table, err := readAt(r, int64(entriesLBA)*sectorSize, int(numEntries*entrySize))
	if err != nil {
		return nil
	}
	parts := []partition{}
	for i := uint32(0); i < numEntries; i++ {
		e := table[i*entrySize : (i+1)*entrySize]
		if bytes.Equal(e[0:16], make([]byte, 16)) {
			continue // unused entry
		}
		parts = append(parts, partition{
			Number: int(i) + 1,
			UUID:   guidString(e[16:32]),
			Label:  utf16String(e[56:128]),
		})
	}
	return parts
}

func probeMbr(r io.ReaderAt) []partition {
	mbr, err := readAt(r, 0, 512)
	if err != nil || mbr[510] != 0x55 || mbr[511] != 0xAA {
		return nil
	}
	diskID := binary.LittleEndian.Uint32(mbr[440:])
	if diskID == 0 {
		return nil
	}
	parts := []partition{}
	for i := 0; i < 4; i++ {
		e := mbr[446+i*16 : 446+(i+1)*16]
		if e[4] == 0 {
			continue // empty entry
		}
		parts = append(parts, partition{
			Number: i + 1,
			UUID:   fmt.Sprintf("%08x-%02d", diskID, i+1),
		})
	}
	return parts
}

// partitionDevice is the device name of partition number n of disk name.
func partitionDevice(disk string, n int) string {
	if len(disk) > 0 && disk[len(disk)-1] >= '0' && disk[len(disk)-1] <= '9' {
		return disk + "p" + strconv.Itoa(n)
	}
	return disk + strconv.Itoa(n)
}

// blockDevices lists the block devices (disks and partitions) of the system
// in the order of /proc/partitions.
func blockDevices() []string {
	data, err := os.ReadFile("/proc/partitions")
	if err != nil {
		return nil
	}
	devices := []string{}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 4 || fields[0] == "major" {
			continue
		}
		devices = append(devices, fields[3])
	}
	return devices
}

func canonicalPath(path string) string {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved
	}
	return path
}

// evaluateTag resolves a "TAG=value" spec to a device path, like
// blkid_evaluate_tag. The udev symlinks are tried first, then the block
// devices of the system are probed.
func evaluateTag(tag, value string) string {
	tag = strings.ToUpper(tag)
	byDir := map[string]string{
		"LABEL":     "by-label",
		"UUID":      "by-uuid",
		"PARTLABEL": "by-partlabel",
		"PARTUUID":  "by-partuuid",
	}
	dir, ok := byDir[tag]
	if !ok {
		return ""
	}
	link := filepath.Join("/dev/disk", dir, udevEncode(value))
	if _, err := os.Stat(link); err == nil {
		return canonicalPath(link)
	}

	for _, name := range blockDevices() {
		dev := "/dev/" + name
		f, err := os.Open(dev)
		if err != nil {
			continue
		}
		var found string
		switch tag {
		case "LABEL", "UUID":
			if p := probeFs(f); p != nil {
				if (tag == "LABEL" && p.Label == value) || (tag == "UUID" && strings.EqualFold(p.UUID, value)) {
					found = dev
				}
			}
		case "PARTLABEL", "PARTUUID":
			for _, part := range probePartitions(f) {
				if (tag == "PARTLABEL" && part.Label == value) || (tag == "PARTUUID" && strings.EqualFold(part.UUID, value)) {
					found = "/dev/" + partitionDevice(name, part.Number)
					break
				}
			}
		}
		f.Close()
		if found != "" {
			if _, err := os.Stat(found); err == nil {
				return found
			}
		}
	}
	return ""
}

// udevEncode escapes a label the way udev does for /dev/disk/by-label
// symlinks: everything but letters, digits, '#', '+', '-', '.', ':', '=',
// '@' and '_' becomes \xNN.
func udevEncode(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9',
			c == '#', c == '+', c == '-', c == '.', c == ':', c == '=', c == '@', c == '_':
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, "\\x%02x", c)
		}
	}
	return b.String()
}

func resolveDeviceSpec(spec string) string {
	if i := strings.IndexByte(spec, '='); i > 0 {
		return evaluateTag(spec[:i], spec[i+1:])
	}
	if spec == "" {
		return ""
	}
	return canonicalPath(spec)
}

func getFsTypeProbe(device string) (string, error) {
	f, err := os.Open(device)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if p := probeFs(f); p != nil {
		return p.Type, nil
	}
	return "", fmt.Errorf("Error while getting fstype")
}
