package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
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

func (is *IfStatement) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {
	return nil
}