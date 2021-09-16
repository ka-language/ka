package ast

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Function struct {
	Params  []*VarDecl //parameter list
	RetType *ASTNode   //return type
	Body    *Block     //function body
}

func (f *Function) Parse(lex []tokenizer.Token, i *int) (e error) {
	*i++ //skip the "fn" token

	//read the return type
	//fn int() {}
	//will also work, because if no braces are present, the next token is returned, and the brace matcher exits
	//if the next value is a variable name, then we know it's a void return type
	//so we will skip the return type

	if lex[*i].Type != "(" {
		rt, e := groupsToAST(groupSpecific(lex, i, []string{"("}))
		f.RetType = rt[0]
		if e != nil {
			return e
		}
	}

	if lex[*i].Type != "(" { //it has to be a parenthesis for the paramlist
		return errors.New("functions require a parameter list")
	}

	p, e := groupsToAST(grouper(braceMatcher(lex, i, []string{"("}, []string{")"}, false, "")))
	sub := p[0].Group.(*Block).Sub
	plist := make([]*VarDecl, len(sub))

	for k, v := range sub {

		switch g := v.Group.(type) {
		case *Operation:
			if g.OpType != ":" {
				return errors.New("invalid syntax: named parameters must have a type")
			}

			plist[k] = &VarDecl{
				Name: v.Left[0].Group.(*VarRef).Name,
				Type: v.Right[0],
			}
		default:
			plist[k] = &VarDecl{
				Type: v,
			}

		}
	}

	f.Params = plist

	if e != nil {
		return e
	}

	*i++

	if lex[*i].Type != "{" {
		*i-- //move back because there is no brace
		return nil
	}

	f.Body = grouper(braceMatcher(lex, i, []string{"{"}, []string{"}"}, false, ""))[0].(*Block)

	if e != nil {
		return e
	}

	return nil
}

func (f *Function) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {
	var rt types.Type = types.Void //defaults to void

	if f.RetType != nil {
		rt = f.RetType.Group.Compile(compiler, class, f.RetType, block).Type()
	}

	var params = make([]*ir.Param, len(f.Params))

	for k, v := range f.Params {
		params[k] = ir.NewParam(
			v.Name,
			v.Type.Group.Compile(compiler, class, v.Type, block).Type(),
		)
	}

	rf := compiler.Module.NewFunc(compiler.TmpVar(), rt, params...)

	if f.Body != nil {
		fblock := rf.NewBlock("")
		f.Body.Compile(compiler, class, nil, fblock)

		//if there is no return type (void) append a `return void`
		if f.RetType == nil {
			fblock.NewRet(nil)
		}
	}

	//if no body was provided, the function was being used as a type
	return data.NewFunc(rf)
}
