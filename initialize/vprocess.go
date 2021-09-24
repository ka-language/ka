package initialize

import (
	"github.com/tusklang/tusk/ast"
)

//decl structure used to store variable declarations
type decl struct {
	nname string //new name
}

//util function to merge two varmaps
func mergemap(m1, m2 map[string]decl) (fin map[string]decl) {

	fin = make(map[string]decl)

	for k, v := range m1 {
		fin[k] = v
	}

	for k, v := range m2 {
		fin[k] = v
	}

	return
}

func (p *VarProcessor) process(tree []*ast.ASTNode, declared map[string]decl) {

	var curscope = make(map[string]decl)

	for _, v := range tree {

		switch g := v.Group.(type) {
		case *ast.VarDecl:

			m := mergemap(declared, curscope)

			if _, exists := m[g.Name]; exists {
				//error
				//variable with that name has already been declared
			}

			p.process([]*ast.ASTNode{g.Type}, m)
			p.process([]*ast.ASTNode{g.Value}, m)

			nname := p.nextvar()
			curscope[g.Name] = decl{
				nname: nname,
			}
			g.Name = nname
		case *ast.VarRef:

			//check both the outer declarations and current scope for the variable reference
			d, ex1 := declared[g.Name]
			cs, ex2 := curscope[g.Name]

			if !(ex1 || ex2) {
				//error
				//there isn't a variable declared with that name
			}

			//if the outer scope doesn't include the var ref, it's in the current scope
			if !ex1 {
				d = cs
			}

			g.Name = d.nname //rename the variable in the ast
		case *ast.Block:
			p.process(g.Sub, mergemap(declared, curscope))
		case *ast.Function:
			if g.Body != nil {

				//do stuff with function parameters
				paramMap := make(map[string]decl)

				for k, v := range g.Params {
					newName := p.nextvar()
					paramMap[v.Name] = decl{
						nname: newName,
					}
					g.Params[k].Name = newName
				}
				///////////////////////////////////

				p.process(g.Body.Sub, mergemap(mergemap(declared, curscope), paramMap))
			}
		case *ast.Operation:
			m := mergemap(declared, curscope)
			p.process(v.Left, m)

			//if it's the dot operator, only check the left side
			if g.OpType != "." {
				p.process(v.Right, m)
			}
		}
	}
}

func (p *VarProcessor) ProcessVars(file *File) {

	var globals = make(map[string]decl)

	for _, v := range file.Globals {
		//add all the globals

		//only add them if they're static
		//because instance variables are accesses with self.<varname>
		if v.IsStatic {
			globals[v.Value.Name] = decl{
				nname: v.Value.Name,
			}
		}

	}

	for _, v := range file.Globals {
		p.process([]*ast.ASTNode{v.Value.Value}, mergemap(p.predecl, globals)) //process the declaration's assigned value
	}

}
