package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type IfStatement struct {
	Condition []*ASTNode
	Body      []*ASTNode
}

func (is *IfStatement) Parse(lex []tokenizer.Token, i *int) error {
	return ifwhileParse(is, lex, i)
}

func (is *IfStatement) SetCond(g []*ASTNode) {
	is.Condition = g
}

func (is *IfStatement) SetBody(g []*ASTNode) {
	is.Body = g
}

func (is *IfStatement) Type() string {
	return "if"
}
