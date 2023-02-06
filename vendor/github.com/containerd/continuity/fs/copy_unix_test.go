//go:build linux || darwin
// +build linux darwin

/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package fs

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/containerd/continuity/fs/fstest"
	"github.com/containerd/continuity/sysx"
	"golang.org/x/sys/unix"
)

func assertXAttr(t *testing.T, dir, xattr, xval string, xerr error) {
	t.Helper()
	value, err := sysx.Getxattr(dir, xattr)
	switch {
	case xerr == nil && err != nil:
		t.Errorf("err=%v; expected val=%s", err, xval)
	case xerr == nil && err == nil && string(value) != xval:
		t.Errorf("val=%s; expected val=%s", value, xval)
	case xerr != nil && err != xerr:
		t.Errorf("val=%s, err=%v; expected err=%v", value, err, xerr)
	}
}

func TestCopyDirWithXAttrExcludes(t *testing.T) {
	src := t.TempDir()
	if err := fstest.Apply(
		fstest.SetXAttr(".", "user.test-1", "one"),
		fstest.SetXAttr(".", "user.test-2", "two"),
		fstest.SetXAttr(".", "user.test-x", "three-four-five"),
	).Apply(src); err != nil {
		t.Fatal(err)
	}

	t.Run("none", func(t *testing.T) {
		dst := t.TempDir()
		err := CopyDir(dst, src, WithXAttrExclude())
		if err != nil {
			t.Fatal(err)
		}
		assertXAttr(t, dst, "user.test-1", "one", nil)
		assertXAttr(t, dst, "user.test-2", "two", nil)
		assertXAttr(t, dst, "user.test-x", "three-four-five", nil)
	})

	t.Run("some", func(t *testing.T) {
		dst := t.TempDir()
		err := CopyDir(dst, src, WithXAttrExclude("user.test-x"))
		if err != nil {
			t.Fatal(err)
		}
		assertXAttr(t, dst, "user.test-1", "one", nil)
		assertXAttr(t, dst, "user.test-2", "two", nil)
		assertXAttr(t, dst, "user.test-x", "", sysx.ENODATA)
	})
}

func TestCopyIrregular(t *testing.T) {
	var prepared int
	prepareSrc := func(src string) {
		f0Pipe := filepath.Join(src, "f0.pipe")
		if err := unix.Mkfifo(f0Pipe, 0600); err != nil {
			t.Fatal(err)
		}
		prepared++
		f1Normal := filepath.Join(src, "f1.normal")
		if err := os.WriteFile(f1Normal, []byte("content of f1.normal"), 0600); err != nil {
			t.Fatal(err)
		}
		prepared++
		f2Socket := filepath.Join(src, "f2.sock")
		if err := unix.Mknod(f2Socket, 0600|unix.S_IFSOCK, 0); err != nil {
			t.Fatal(err)
		}
		prepared++
		f3Dev := filepath.Join(src, "f3.dev")
		if err := unix.Mknod(f3Dev, 0600|unix.S_IFCHR, 42); err != nil {
			t.Logf("skipping testing S_IFCHR: %v", err)
		} else {
			prepared++
		}
	}

	verifyDst := func(dst string) {
		entries, err := os.ReadDir(dst)
		if err != nil {
			t.Fatal(err)
		}
		var verified int
		for _, f := range entries {
			name := f.Name()
			full := filepath.Join(dst, name)
			fi, err := os.Stat(full)
			if err != nil {
				t.Fatal(err)
			}
			mode := fi.Mode()
			switch name {
			case "f0.pipe":
				if mode&os.ModeNamedPipe != os.ModeNamedPipe {
					t.Fatalf("unexpected mode of %s: %v", name, mode)
				}
			case "f1.normal":
				b, err := os.ReadFile(full)
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != "content of f1.normal" {
					t.Fatalf("unexpected content of %s: %q", name, string(b))
				}
			case "f2.sock":
				if mode&os.ModeSocket != os.ModeSocket {
					t.Fatalf("unexpected mode of %s: %v", name, mode)
				}
			case "f3.dev":
				if mode&os.ModeDevice != os.ModeDevice {
					t.Fatalf("unexpected mode of %s: %v", name, mode)
				}
				sys, ok := fi.Sys().(*syscall.Stat_t)
				if !ok {
					t.Fatalf("unexpected type: %v", fi.Sys())
				}
				if sys.Rdev != 42 {
					t.Fatalf("unexpected rdev of %s: %d", name, sys.Rdev)
				}
			}
			verified++
		}
		if verified != prepared {
			t.Fatalf("prepared %d files, verified %d files", prepared, verified)
		}
	}

	src := t.TempDir()
	dst := t.TempDir()
	prepareSrc(src)
	if err := CopyDir(dst, src); err != nil {
		t.Fatal(err)
	}
	verifyDst(dst)
}
