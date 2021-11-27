package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type VarRef struct {
	Name string
	tok  tokenizer.Token
}

func (vr *VarRef) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {
	vr.Name = lex[*i].Name
	vr.tok = lex[*i]
	return nil
}

func (vr *VarRef) GetMTok() tokenizer.Token {
	return vr.tok
}

func (vr *VarRef) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	fetched := compiler.FetchVar(vr.Name)

	if fetched == nil {

		//check the class' static variables if there is no variable declared with x name

		_fetched := class.Static[vr.Name]

		if _fetched != nil {
			fetched = _fetched.Value
		}

		if fetched == nil {
			//if there still isn't a variable with that name, it's an "undeclared variable"
			return data.NewUndeclaredVar(vr.Name)
		}

	}

	return fetched
}
