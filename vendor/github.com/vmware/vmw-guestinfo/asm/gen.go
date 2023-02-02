// Copyright 2020 Google, Inc. All Rights Reserved.
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

package asm

import (
	"flag"
	"go/types"
	"log"
	"os"

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/buildtags"
	"github.com/mmcloughlin/avo/gotypes"
	"github.com/mmcloughlin/avo/ir"
)

// CLI utilities for ASM generator

// AvoContext is an extended context for Avo that conveys information about
// the target architecture. In particular this allows some level of
// abstraction over registers.
type AvoContext struct {
	*build.Context
	Arch string
}

// Missing instructions from avo.

// INL is the INL instruction.
func (ctx *AvoContext) INL() {
	ctx.Instruction(&ir.Instruction{Opcode: "INL"})
}

// INSB is the INSB instruction.
func (ctx *AvoContext) INSB() {
	ctx.Instruction(&ir.Instruction{Opcode: "INSB"})
}

// OUTSB is the OUTSB instruction.
func (ctx *AvoContext) OUTSB() {
	ctx.Instruction(&ir.Instruction{Opcode: "OUTSB"})
}

// REP is the REP instruction.
func (ctx *AvoContext) REP() {
	ctx.Instruction(&ir.Instruction{Opcode: "REP"})
}

// AvoMainFunc is the type of functions that can be used with GenAsm.
type AvoMainFunc func(ctx *AvoContext) error

// GenAsm is our main entry point for ASM code generation.
// It augments avo flags with a "-arch", the content of which is
// exposed in AvoContext. That value also drives the alignments in
// the generated code.
func GenAsm(f AvoMainFunc) {
	var arch string

	// reset the flagset to allow NewFlags to work.
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var flags = build.NewFlags(flag.CommandLine)
	flag.CommandLine.StringVar(&arch, "arch", "amd64", "arch to use")

	flag.Parse()

	// This is needed to get proper alignment of the return values.
	// amd64 is the global default in avo, but not modifying it might
	// result in misalignment of the stack.
	gotypes.Sizes = types.SizesFor("gc", arch)

	cfg := flags.Config()

	ctx := &AvoContext{
		Context: build.NewContext(),
		Arch:    arch,
	}

	ctx.Constraint(buildtags.Opt("gc"))

	if err := f(ctx); err != nil {
		log.Fatal(err)
	}

	os.Exit(build.Main(cfg, ctx.Context))
}
