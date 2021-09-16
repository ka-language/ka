package initialize

import (
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
)

type GlobalDecl struct {

	/*
		0: public
		1: protected
		2: private
	*/
	Access int

	IsStatic bool //if the global is a static or instance
	Value    *ast.VarDecl
}

type File struct {
	Name        string
	Globals     []GlobalDecl
	StructType  *types.StructType
	Constructor *ast.Construct
}
