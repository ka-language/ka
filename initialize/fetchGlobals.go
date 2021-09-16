package initialize

import (
	"github.com/tusklang/tusk/ast"
)

func fetchGlobals(tree []*ast.ASTNode, file *File, access int, isStatic bool) {

	for _, v := range tree {

		switch g := v.Group.(type) {
		case *ast.VarDecl:
			file.Globals = append(file.Globals, GlobalDecl{
				Access:   access,
				IsStatic: isStatic,
				Value:    g,
			})
		case *ast.Public:
			fetchGlobals([]*ast.ASTNode{g.Declaration}, file, 0, isStatic)
		case *ast.Protected:
			fetchGlobals([]*ast.ASTNode{g.Declaration}, file, 1, isStatic)
		case *ast.Static:
			fetchGlobals([]*ast.ASTNode{g.Declaration}, file, access, true)
		case *ast.Construct:
			file.Constructor = g
		}
	}
}
