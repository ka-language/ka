package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type Public struct {
	Declaration *ASTNode
}

func (p *Public) Parse(lex []tokenizer.Token, i *int) (e error) {
	return parseAccessSpec(p, lex, i)
}

func (p *Public) SetDecl(node *ASTNode) {
	p.Declaration = node
}
