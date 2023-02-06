//go:build !darwin && !windows
// +build !darwin,!windows

package mount

import (
	"errors"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/moby/sys/mountinfo"
	"golang.org/x/sys/unix"
)

func TestMountOptionsParsing(t *testing.T) {
	options := "noatime,ro,noexec,size=10k"

	flag, data := parseOptions(options)

	if data != "size=10k" {
		t.Fatalf("Expected size=10 got %s", data)
	}

	expected := NOATIME | RDONLY | NOEXEC

	if flag != expected {
		t.Fatalf("Expected %d got %d", expected, flag)
	}
}

func TestMounted(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("root required")
	}

	var (
		tmp        = t.TempDir()
		sourceDir  = path.Join(tmp, "source")
		targetDir  = path.Join(tmp, "target")
		sourcePath = path.Join(sourceDir, "file.txt")
		targetPath = path.Join(targetDir, "file.txt")
	)

	if err := os.Mkdir(sourceDir, 0o777); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(targetDir, 0o777); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(sourcePath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(targetPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Mount(sourceDir, targetDir, "none", "bind,rw"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Unmount(targetDir); err != nil {
			t.Fatal(err)
		}
	}()

	mounted, err := mountinfo.Mounted(targetDir)
	if err != nil {
		t.Fatal(err)
	}
	if !mounted {
		t.Fatalf("Expected %s to be mounted", targetDir)
	}
	if _, err := os.Stat(targetDir); err != nil {
		t.Fatal(err)
	}
}

func TestMountTmpfsOptions(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("root required")
	}

	testCases := []struct {
		opts       string
		expected   string
		unexpected string
	}{
		{
			opts:       "exec",
			unexpected: "noexec",
		},
		{
			opts:       "noexec",
			expected:   "noexec",
			unexpected: "exec",
		},
	}

	target := t.TempDir()
	for _, tc := range testCases {
		t.Run(tc.opts, func(t *testing.T) {
			if err := Mount("tmpfs", target, "tmpfs", tc.opts); err != nil {
				t.Fatal(err)
			}
			defer ensureUnmount(t, target)

			mounts, err := mountinfo.GetMounts(mountinfo.SingleEntryFilter(target))
			if err != nil {
				t.Fatal(err)
			}
			if len(mounts) != 1 {
				t.Fatal("Mount point ", target, " not found")
			}
			entry := mounts[0]
			opts := "," + entry.Options + ","
			if tc.expected != "" && !strings.Contains(opts, ","+tc.expected+",") {
				t.Fatal("Expected option ", tc.expected, " missing from ", entry.Options)
			}
			if tc.unexpected != "" && strings.Contains(opts, ","+tc.unexpected+",") {
				t.Fatal("Unexpected option ", tc.unexpected, " in ", entry.Options)
			}
		})
	}
}

func TestMountReadonly(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("root required")
	}

	var (
		tmp        = t.TempDir()
		sourceDir  = path.Join(tmp, "source")
		targetDir  = path.Join(tmp, "target")
		sourcePath = path.Join(sourceDir, "file.txt")
		targetPath = path.Join(targetDir, "file.txt")
	)

	if err := os.Mkdir(sourceDir, 0o777); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(targetDir, 0o777); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(sourcePath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(targetPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Mount(sourceDir, targetDir, "none", "bind,ro"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Unmount(targetDir); err != nil {
			t.Fatal(err)
		}
	}()

	if err := os.WriteFile(targetPath, []byte("hello"), 0o644); err == nil {
		t.Fatal("Should not be able to open a ro file as rw")
	}
}

func TestMergeTmpfsOptions(t *testing.T) {
	options := []string{"noatime", "ro", "size=10k", "defaults", "noexec", "atime", "defaults", "rw", "rprivate", "size=1024k", "slave", "exec"}
	expected := []string{"atime", "rw", "size=1024k", "slave", "exec"}
	merged, err := MergeTmpfsOptions(options)
	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(merged) {
		t.Fatalf("Expected %s got %s", expected, merged)
	}
	for index := range merged {
		if merged[index] != expected[index] {
			t.Fatalf("Expected %s for the %dth option, got %s", expected, index, merged)
		}
	}

	options = []string{"noatime", "ro", "size=10k", "atime", "rw", "rprivate", "size=1024k", "slave", "size", "exec"}
	_, err = MergeTmpfsOptions(options)
	if err == nil {
		t.Fatal("Expected error got nil")
	}
}

func TestRecursiveUnmountTooGreedy(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("root required")
	}

	tmp := t.TempDir()
	// Create a bunch of tmpfs mounts. Make sure "dir" itself is not
	// a mount point, or we'll hit the fast path in RecursiveUnmount.
	dirs := []string{"dir-other", "dir/subdir1", "dir/subdir1/subsub", "dir/subdir2/subsub"}
	for _, d := range dirs {
		dir := path.Join(tmp, d)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := Mount("tmpfs", dir, "tmpfs", ""); err != nil {
			t.Fatal(err)
		}
		//nolint:errcheck
		defer Unmount(dir)
	}
	// sanity check
	mounted, err := mountinfo.Mounted(path.Join(tmp, "dir-other"))
	if err != nil {
		t.Fatalf("[pre-check] error from mountinfo.mounted: %v", err)
	}
	if !mounted {
		t.Fatal("[pre-check] expected dir-other to be mounted, but it's not")
	}
	// Unmount dir, make sure dir-other is still mounted.
	if err := RecursiveUnmount(path.Join(tmp, "dir")); err != nil {
		t.Fatal(err)
	}
	mounted, err = mountinfo.Mounted(path.Join(tmp, "dir-other"))
	if err != nil {
		t.Fatalf("error from mountinfo.mounted: %v", err)
	}
	if !mounted {
		t.Fatal("expected dir-other to be mounted, but it's not")
	}
}

func TestRecursiveUnmount_SubMountFailsToUnmount(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("root required")
	}

	var (
		tmp        = t.TempDir()
		parent     = tmp + "/sub1"
		child      = tmp + "/sub1/sub2"
		grandChild = tmp + "/sub1/sub2/sub3"
	)

	err := os.MkdirAll(grandChild, 0o700)
	if err != nil {
		t.Fatal(err)
	}

	// Create a set of mounts that should result in RecursiveUnmount failure,
	// caused by the fact that the grandchild mount is shadowed by the child mount,
	// and the child mount is shadowed by the parent mount. So. these two mounts
	// are listed in mountinfo, but since they are unreachable, unmount will fail.
	toMount := []string{grandChild, child, parent}
	for _, dir := range toMount {
		dir := dir
		if err := Mount("tmpfs", dir, "tmpfs", ""); err != nil {
			t.Fatal(err)
		}
		defer Unmount(dir) //nolint:errcheck
	}

	// unmount shadowed mounts
	shadowedMounts := []string{child, grandChild}
	for _, shadowedMount := range shadowedMounts {
		t.Run(shadowedMount, func(t *testing.T) {
			err := RecursiveUnmount(shadowedMount)
			if !errors.Is(err, unix.ENOENT) {
				t.Fatalf("expected submount(shadowed) %s to return unix.ENOENT, got %v", shadowedMount, err)
			}
		})
	}
}
