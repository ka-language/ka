package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type NullValue struct{}

func (nv *NullValue) Parse(lex []tokenizer.Token, i *int) error {
	return nil
}
