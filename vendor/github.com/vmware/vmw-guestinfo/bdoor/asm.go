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
	"fmt"
	"strings"

	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/reg"

	"github.com/vmware/vmw-guestinfo/asm"
)

// Doc of the golang plan9 assembler
// http://p9.nyx.link/labs/sys/doc/asm.html
//
// A good primer of how to write golang with some plan9 flavored assembly
// http://www.doxsey.net/blog/go-and-assembly
//
// Some x86 references
// http://www.eecg.toronto.edu/~amza/www.mindsec.com/files/x86regs.html
// https://cseweb.ucsd.edu/classes/sp10/cse141/pdf/02/S01_x86_64.key.pdf
// https://en.wikibooks.org/wiki/X86_Assembly/Other_Instructions
//
// (This one is invaluable.  Has a working example of how a standard function
// call looks on the stack with the associated assembly.)
// https://www.recurse.com/blog/7-understanding-c-by-learning-assembly
//
// Reference with raw form of the Opcode
// http://x86.renejeschke.de/html/file_module_x86_id_139.html
//
// Massive x86_64 reference
// http://ref.x86asm.net/coder64.html#xED

type config struct {
	wordSize  uint8
	registers map[string]reg.Register
}

func setup(ctx *asm.AvoContext) (*config, error) {
	switch ctx.Arch {
	case "amd64":
		return &config{
			wordSize: 64,
			registers: map[string]reg.Register{
				"ax": reg.RAX,
				"bx": reg.RBX,
				"cx": reg.RCX,
				"dx": reg.RDX,
				"si": reg.RSI,
				"di": reg.RDI,
				"bp": reg.RBP,
			},
		}, nil
	case "386":
		return &config{
			wordSize: 32,
			registers: map[string]reg.Register{
				"ax": reg.EAX,
				"bx": reg.EBX,
				"cx": reg.ECX,
				"dx": reg.EDX,
				"si": reg.ESI,
				"di": reg.EDI,
				"bp": reg.EBP,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported architecture: %s", ctx.Arch)
	}
}

func code(ctx *asm.AvoContext) error {
	cfg, err := setup(ctx)
	if err != nil {
		return err
	}

	// for all the following functions, we use those 7 registers as both
	// input and output.
	regs := []string{"ax", "bx", "cx", "dx", "si", "di", "bp"}

	funcs := []struct {
		name string
		body func()
	}{
		{"bdoor_inout", func() {
			ctx.Comment("IN to DX from AX")
			ctx.INL()
		}},
		{"bdoor_hbout", func() {
			ctx.CLD()
			ctx.REP()
			ctx.OUTSB()
		}},
		{"bdoor_hbin", func() {
			ctx.CLD()
			ctx.REP()
			ctx.INSB()
		}},
		{"bdoor_inout_test", func() {}},
	}

	for _, f := range funcs {
		// generate function signature. By convention we use:
		// - $REG as input
		// - and ret$REG as output.
		proto := fmt.Sprintf(
			"func(%s uint%d) (%s uint%d)",
			strings.Join(regs, ","),
			cfg.wordSize,
			"ret"+strings.Join(regs, ",ret"),
			cfg.wordSize,
		)

		ctx.Function(f.name)
		ctx.Attributes(NOSPLIT | WRAPPER)
		ctx.SignatureExpr(proto)

		// load all registers
		for _, r := range regs {
			ctx.Load(ctx.Param(r), cfg.registers[r])
		}

		// execute backdoor function
		f.body()

		// store all registers
		for _, r := range regs {
			ctx.Store(cfg.registers[r], ctx.Return("ret"+r))
		}

		ctx.RET()
	}

	return nil
}

func main() {
	asm.GenAsm(code)
}
