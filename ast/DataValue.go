package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type DataValue struct {
	Value tokenizer.Token
}

func (dv *DataValue) Parse(lex []tokenizer.Token, i *int) error {
	dv.Value = lex[*i]
	return nil
}
