// Copyright 2016 VMware, Inc. All Rights Reserved.
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

package bdoor

import (
	"testing"

	"github.com/vmware/vmw-guestinfo/internal"
)

func TestBdoorArgAlignment(t *testing.T) {
	a := uint64(0xFFFFFFFF0000022)
	b := uint64(33)
	c := uint64(44)
	d := uint64(55)
	si := uint64(0xFFFFFFFF0000066)
	di := uint64(0xFFFAAFFF0000077)
	bp := uint64(0xFFFFFFFFAAAAAAA)

	oa, ob, oc, od, osi, odi, obp := bdoor_inout_test(a, b, c, d, si, di, bp)

	internal.AssertEqual(t, a, oa)
	internal.AssertEqual(t, b, ob)
	internal.AssertEqual(t, c, oc)
	internal.AssertEqual(t, d, od)
	internal.AssertEqual(t, si, osi)
	internal.AssertEqual(t, di, odi)
	internal.AssertEqual(t, bp, obp)
}
