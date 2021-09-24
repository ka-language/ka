package ast

import (
	"github.com/llir/llvm/ir"
)

type Compiler struct {
	Module   *ir.Module //llvm module
	OS, ARCH string     //operating system and architecture
}
