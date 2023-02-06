// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vmcheck

import (
	"encoding/binary"
	"fmt"
	"testing"
)

var (
	intelID = []byte("GenuineIntel")
	vmwID   = []byte("VMwareVMware")

	errPerm   = fmt.Errorf("operation not permitted")
	errAccess = fmt.Errorf("access denied")
)

func TestVirtualWorld(t *testing.T) {
	leUint32 := binary.LittleEndian.Uint32
	data := []struct {
		scenario   string
		hv         bool
		id         []byte
		portAccess error
		bdoorKnock error
		resp       bool
	}{
		{"vmware", true, vmwID, nil, nil, true},
		{"vmware no backdoor", true, vmwID, errPerm, errAccess, false},
		{"physical", false, intelID, errPerm, errAccess, false},
	}

	for _, tt := range data {
		tt := tt
		t.Run(tt.scenario, func(t *testing.T) {
			p := platform{
				accessPorts: func() error { return tt.portAccess },
				cpuid: func(eax, _ uint32) (uint32, uint32, uint32, uint32) {
					if eax == 0x1 {
						if tt.hv {
							return 0, 0, uint32(1 << 31), 0
						}
						return 0, 0, 0, 0
					}
					if eax == 0x40000000 {
						return 0, leUint32(tt.id[:]), leUint32(tt.id[4:]), leUint32(tt.id[8:])
					}
					t.Fatal(fmt.Errorf("unexpected cpuid call"))
					return 0, 0, 0, 0
				},
				knock: func() (bool, error) {
					if tt.bdoorKnock == nil {
						return true, nil
					}
					return false, tt.bdoorKnock
				},
			}
			v, err := p.isVirtualWorld(false)
			if tt.resp {
				if !v || err != nil {
					t.Fatal("expected to be virtual")
				}
				return
			}

			if v {
				t.Fatal("wrongly detected as virtual")
			}

			if err != nil {
				if tt.portAccess != nil && err == tt.portAccess {
					return
				}
				if tt.bdoorKnock != nil && err == tt.bdoorKnock {
					return
				}
				t.Fatal("unexpected error")
			}
		})
	}
}
