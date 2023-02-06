// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
// Copyright 2020 Google LLC
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

// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/reg"

	"github.com/vmware/vmw-guestinfo/asm"
)

func code(ctx *asm.AvoContext) error {
	ctx.Function("cpuid_low")
	ctx.Attributes(NOSPLIT)
	ctx.SignatureExpr("func(arg1, arg2 uint32) (eax, ebx, ecx, edx uint32)")

	ctx.Comment("From https://github.com/intel-go/cpuid/blob/master/cpuidlow_amd64.s")

	ctx.Load(ctx.Param("arg1"), reg.EAX)
	ctx.Load(ctx.Param("arg2"), reg.ECX)

	ctx.CPUID()

	ctx.Store(reg.EAX, ctx.Return("eax"))
	ctx.Store(reg.EBX, ctx.Return("ebx"))
	ctx.Store(reg.ECX, ctx.Return("ecx"))
	ctx.Store(reg.EDX, ctx.Return("edx"))

	ctx.RET()

	return nil
}

func main() {
	asm.GenAsm(code)
}
