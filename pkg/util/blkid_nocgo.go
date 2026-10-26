//go:build linux && !cgo
// +build linux,!cgo

package util

// ResolveDevice resolves a device spec (a path or LABEL=, UUID=, PARTLABEL=,
// PARTUUID=) to a device path without libblkid, see blkid_probe.go.
func ResolveDevice(spec string) string {
	return resolveDeviceSpec(spec)
}

// GetFsType returns the file system type found on device without libblkid,
// see blkid_probe.go.
func GetFsType(device string) (string, error) {
	return getFsTypeProbe(device)
}
