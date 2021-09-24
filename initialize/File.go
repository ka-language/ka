package initialize

import (
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
	Constructor *ast.Construct

	pkg *Directory //parent package
}

func (f *File) Parent() FSObj {
	return f.pkg
}
