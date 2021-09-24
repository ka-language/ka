package ast

import (
	"errors"

	"github.com/tusklang/tusk/tokenizer"
)

type Construct struct {
	FnObj *Function
}

func (c *Construct) Parse(lex []tokenizer.Token, i *int) error {

	var fnobj = &Function{}
	e := fnobj.Parse(lex, i) //functions and constructors are (surprisingly enough :p) structured the same

	if e != nil { //if the function parse returned an error
		return e
	}

	if fnobj.RetType != nil { //constructors cannot return anything
		return errors.New("constructor cannot include a return type")
	}

	c.FnObj = fnobj

	return nil
}
