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
	"github.com/vmware/vmw-guestinfo/internal"
	"testing"
)

func TestSetWord(t *testing.T) {
	inLow := uint16(0xEEFF)
	inHigh := uint16(0xBBBB)

	out := &UInt32{}
	//out.SetWord(uint32(0xBBBBEEFF))
	out.Low = inLow
	out.High = inHigh

	internal.AssertEqual(t, inLow, out.Low)
	internal.AssertEqual(t, inHigh, out.High)

	internal.AssertEqual(t, uint32(0xBBBBEEFF), out.Word())
}

func TestQuadToHighLow(t *testing.T) {
	in := uint64(0xFFFFFFFF0000000A)

	var u UInt64
	u.SetQuad(in)
	internal.AssertEqual(t, uint32(in), u.Low.Word())
	internal.AssertEqual(t, uint32(in>>32), u.High.Word())
	internal.AssertEqual(t, in, u.Quad())
}

func TestHighLowToQuad(t *testing.T) {
	inHigh := uint16(0xff)
	inLow := uint16(0xaa)

	u := UInt64{
		High: UInt32{High: inHigh},
		Low:  UInt32{Low: inLow},
	}

	internal.AssertEqual(t, (uint64(inHigh)<<48)+uint64(inLow), u.Quad())
}
