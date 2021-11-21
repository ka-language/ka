package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Array struct {
	siz *ASTNode
	typ *ASTNode
	arr []*ASTNode

	//if the array is used as an index (var a = [2]i32{1, 2}; a[0];)
	//the first statement is being used an expression
	//the second statement is being used an index
	useAsIndex bool

	//used during compiling
	csiz data.Value
	ctyp data.Type
}

func (a *Array) Parse(lex []tokenizer.Token, i *int) error {
	sizl := braceMatcher(lex, i, []string{"[", "{", "("}, []string{"]", "}", ")"}, true, "")
	*i++
	sizg := grouper(sizl)
	siz, e := groupsToAST(sizg)

	if e != nil {
		return e
	}

	if len(siz) == 1 {
		a.siz = siz[0]
	} else if len(siz) != 0 {
		//error
		//size can't be multiple statements
	}

	if *i >= len(lex) {
		return nil
	}

	if lex[*i].Type == "(" || lex[*i].Type == "varname" {
		typg := groupSpecific(lex, i, nil, 1)
		typ, e := groupsToAST(typg)

		if e != nil {
			return e
		}

		a.typ = typ[0]
	}

	//arrays are written like:
	//var v: []Type = []Type{}
	//if there is no { after the type, then it is being used as the var type
	if lex[*i].Type == "{" {
		arrl := braceMatcher(lex, i, []string{"{"}, []string{"}"}, true, "")
		arrg := grouper(arrl)
		arr, e := groupsToAST(arrg)

		if e != nil {
			return e
		}

		a.arr = arr
	} else {
		//there was no array content
		*i--
	}

	return nil
}

func (a *Array) CompileSlice(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}

func (a *Array) CompileFixedArray(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	sizi := a.csiz.(*data.Integer).GetInt()

	block := function.ActiveBlock
	arrtyp := types.NewArray(uint64(sizi), a.ctyp.Type())
	decl := block.NewAlloca(arrtyp)
	decl.Align = ir.Align(16)

	//fill the array with the values needed
	for k, v := range a.arr {
		vc := v.Group.Compile(compiler, class, v, function)
		gep := block.NewGetElementPtr(arrtyp, decl, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(k)))
		block.NewStore(vc.LLVal(block), gep)
	}

	curlen := block.NewAlloca(types.I32)
	block.NewStore(constant.NewInt(types.I32, sizi), curlen)

	return data.NewFixedArray(a.ctyp, decl, curlen, uint64(sizi))
}

func (a *Array) CompileVariedLengthArray(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}

func (a *Array) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	if a.siz == nil {
		//it's a slice array
		a.ctyp = a.typ.Group.Compile(compiler, class, a.typ, function).(data.Type)
	} else {
		a.csiz = a.siz.Group.Compile(compiler, class, a.siz, function)

		if a.useAsIndex {
			return a.csiz
		}

		a.ctyp = a.typ.Group.Compile(compiler, class, a.typ, function).(data.Type)

		switch a.csiz.(type) {
		case *data.Integer:
			//it's a fixed array
			return a.CompileFixedArray(compiler, class, node, function)
		default:
			//it's a varied length array
		}
	}

	return nil
}
