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
	"context"
	"errors"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/containerd/continuity/fs/fstest"
)

var errNotImplemented = errors.New("check not implemented")

func TestUsage(t *testing.T) {
	align, dirs, err := getTmpAlign(t)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		name string
		fs   fstest.Applier
		size int64
	}
	testCases := []testCase{
		{
			name: "SingleSmallFile",
			fs: fstest.Apply(
				fstest.CreateDir("/dir", 0755),
				fstest.CreateRandomFile("/dir/file", 1, 5, 0644),
			),
			size: dirs(2) + align(5),
		},
		{
			name: "MultipleSmallFile",
			fs: fstest.Apply(
				fstest.CreateDir("/dir", 0755),
				fstest.CreateRandomFile("/dir/file1", 2, 5, 0644),
				fstest.CreateRandomFile("/dir/file2", 3, 5, 0644),
			),
			size: dirs(2) + align(5)*2,
		},
		{
			name: "BiggerFiles",
			fs: fstest.Apply(
				fstest.CreateDir("/dir", 0755),
				fstest.CreateRandomFile("/dir/file1", 4, 5, 0644),
				fstest.CreateRandomFile("/dir/file2", 5, 1024, 0644),
				fstest.CreateRandomFile("/dir/file3", 6, 50*1024, 0644),
			),
			size: dirs(2) + align(5) + align(1024) + align(50*1024),
		},
	}
	if runtime.GOOS != "windows" {
		testCases = append(testCases, []testCase{
			{
				name: "SparseFiles",
				fs: fstest.Apply(
					fstest.CreateDir("/dir", 0755),
					fstest.CreateRandomFile("/dir/file1", 7, 5, 0644),
					createSparseFile("/dir/sparse1", 8, 0644, 5, 1024*1024, 5),
					createSparseFile("/dir/sparse2", 9, 0644, 0, 1024*1024),
					createSparseFile("/dir/sparse2", 10, 0644, 0, 1024*1024*1024, 1024),
				),
				size: dirs(2) + align(5)*3 + align(1024),
			},
			{
				name: "Hardlinks",
				fs: fstest.Apply(
					fstest.CreateDir("/dir", 0755),
					fstest.CreateRandomFile("/dir/file1", 11, 60*1024, 0644),
					fstest.Link("/dir/file1", "/dir/link1"),
				),
				size: dirs(2) + align(60*1024),
			},
			{
				name: "HardlinkSparefile",
				fs: fstest.Apply(
					fstest.CreateDir("/dir", 0755),
					createSparseFile("/dir/file1", 10, 0644, 30*1024, 1024*1024*1024, 30*1024),
					fstest.Link("/dir/file1", "/dir/link1"),
				),
				size: dirs(2) + align(30*1024)*2,
			},
		}...)
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			t1 := t.TempDir()
			if err := tc.fs.Apply(t1); err != nil {
				t.Fatal("Failed to apply base filesystem:", err)
			}

			usage, err := DiskUsage(context.Background(), t1)
			if err != nil {
				t.Fatal(err)
			}
			if usage.Size != tc.size {
				t.Fatalf("Wrong usage size %d, expected %d", usage.Size, tc.size)
			}

			du, err := duCheck(t1)
			if err != nil && err != errNotImplemented {
				t.Fatal("Failed calling du:", err)
			}
			if err == nil && usage.Size != du {
				t.Fatalf("Wrong usage size %d, du reported %d", usage.Size, du)
			}
		})
	}
}

// createSparseFile creates a sparse file filled with random
// bytes for data parts
// The parse alternate data length, hole length, data length, ....
// To start a file as sparse, give an initial data length of 0
func createSparseFile(name string, seed int64, perm os.FileMode, parts ...int64) fstest.Applier {
	return sparseFile{
		name:  name,
		seed:  seed,
		parts: parts,
		perm:  perm,
	}
}

type sparseFile struct {
	name  string
	seed  int64
	parts []int64
	perm  os.FileMode
}

func (sf sparseFile) Apply(root string) (retErr error) {
	fullPath := filepath.Join(root, sf.name)
	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sf.perm)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil && retErr == nil {
			retErr = err
		}
	}()

	rr := rand.New(rand.NewSource(sf.seed))

	parts := sf.parts
	for len(parts) > 0 {
		// Write content
		if parts[0] > 0 {
			_, err = io.Copy(f, io.LimitReader(rr, parts[0]))
			if err != nil {
				return err
			}
		}
		parts = parts[1:]

		if len(parts) > 0 {
			if parts[0] != 0 {
				f.Seek(parts[0], io.SeekCurrent)
			}
			parts = parts[1:]
		}
	}
	return os.Chmod(fullPath, sf.perm)
}
