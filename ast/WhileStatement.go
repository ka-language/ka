package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type WhileStatement struct {
	Condition []*ASTNode
	Body      []*ASTNode
}

func (ws *WhileStatement) Parse(lex []tokenizer.Token, i *int) error {
	return ifwhileParse(ws, lex, i)
}

func (ws *WhileStatement) SetCond(g []*ASTNode) {
	ws.Condition = g
}

func (ws *WhileStatement) SetBody(g []*ASTNode) {
	ws.Body = g
}

func (ws *WhileStatement) Type() string {
	return "while"
}
