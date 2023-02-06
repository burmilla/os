//go:build freebsd || darwin
// +build freebsd darwin

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
	"fmt"
	"os/exec"
)

func duCmd(root string) *exec.Cmd {
	cmd := exec.Command("du", "-s", root)
	/*
		From du(1):
		BLOCKSIZE  If the environment variable BLOCKSIZE is set, and the -h, -k,
		           -m or --si options are not specified, the block counts will be
		           displayed in units of that block size.  If BLOCKSIZE is not
		           set, and the -h, -k, -m or --si options are not specified, the
		           block counts will be displayed in 512-byte blocks.
	*/
	cmd.Env = []string{fmt.Sprintf("BLOCKSIZE=%d", blocksUnitSize)}
	return cmd
}
