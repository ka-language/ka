package ast

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type VarDecl struct {
	Name  string
	Type  *ASTNode
	Value *ASTNode

	declaration *ir.Global
}

func (vd *VarDecl) Parse(lex []tokenizer.Token, i *int) error {

	*i++

	if lex[*i].Type != "varname" {
		return errors.New("expected a variable name")
	}

	vd.Name = lex[*i].Name

	*i++

	//has a specified type
	if lex[*i].Name == ":" {
		*i++
		t, e := groupsToAST(groupSpecific(lex, i, []string{"=", ";"}))
		if e != nil {
			return e
		}
		vd.Type = t[0]
	}

	//has a value assigned to it
	if lex[*i].Name == "=" {
		*i++
		v, e := groupsToAST(grouper(braceMatcher(lex, i, allopeners, allclosers, false, "terminator")))
		vd.Value = v[0]
		if e != nil {
			return e
		}
	}

	*i-- //the outer loop will incremenet for us

	return nil
}

func (vd *VarDecl) getDeclType(compiler *Compiler, class *data.Class, block *ir.Block) data.Type {

	pvtype := vd.Type.Group.Compile(compiler, class, vd.Type, block)

	var vtype data.Type

	switch vt := pvtype.(type) {
	case *data.Class:
		// vtype = data.NewType(vt.Type())
		// vtype.SetTypeName("instance " + vt.TypeString()) //if it's in a variable decl, it's an instance of a class
		vtype = vt
	case *data.Primative:
		vtype = vt
	case *data.Function:
		vtype = vt
	default:
		//error
	}

	return vtype
}

func (vd *VarDecl) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {

	varval := vd.Value.Group.Compile(compiler, class, vd.Value, block)

	vtype := vd.getDeclType(compiler, class, block)

	decl := block.NewAlloca(vtype.Type())

	if llv := varval.LLVal(block); llv != nil {
		block.NewStore(llv, decl)
	}

	dv := data.NewVariable(data.NewInstruction(decl), vtype, false)

	compiler.AddVar(vd.Name, dv)

	return dv
}

func (vd *VarDecl) DeclareGlobal(name string, compiler *Compiler, class *data.Class, static bool) error {

	vtype := vd.getDeclType(compiler, class, compiler.InitBlock)

	vd.declaration = compiler.Module.NewGlobal(name, vtype.Type())
	vd.declaration.Init = vtype.Default()

	nv := data.NewVariable(data.NewInstruction(vd.declaration), vtype, false)

	if static {
		class.Static[vd.Name] = nv
	} else {
		class.Instance.Put(vd.Name, nv)
	}

	return nil
}

//used specifically for global variable declarations
func (vd *VarDecl) CompileGlobal(compiler *Compiler, class *data.Class, block *ir.Block) {
	val := vd.Value.Group.Compile(compiler, class, vd.Value, compiler.InitBlock)
	block.NewStore(val.LLVal(block), vd.declaration)
}
