//go:build !windows
// +build !windows

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
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

func getBsize(root string) (int64, error) {
	var s syscall.Statfs_t
	if err := syscall.Statfs(root, &s); err != nil {
		return 0, err
	}

	return int64(s.Bsize), nil //nolint: unconvert
}

// getTmpAlign returns filesystem specific size alignment functions
// first:  aligns filesize to file usage based on blocks
// second: determines directory usage based on directory count (assumes small directories)
func getTmpAlign(t testing.TB) (func(int64) int64, func(int64) int64, error) {
	t1 := t.TempDir()

	bsize, err := getBsize(t1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get bsize: %w", err)
	}

	align := func(size int64) int64 {
		// Align to blocks
		aligned := (size / bsize) * bsize

		// Add next block if has remainder
		if size%bsize > 0 {
			aligned += bsize
		}

		return aligned
	}

	fi, err := os.Stat(t1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to stat directory: %w", err)
	}

	dirSize := fi.Sys().(*syscall.Stat_t).Blocks * blocksUnitSize
	dirs := func(count int64) int64 {
		return count * dirSize
	}

	return align, dirs, nil
}

func duCheck(root string) (usage int64, err error) {
	cmd := duCmd(root)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	if len(out) == 0 {
		return 0, errors.New("no du output")
	}
	size := strings.Fields(string(out))[0]
	blocks, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return 0, err
	}

	return blocks * blocksUnitSize, nil
}
