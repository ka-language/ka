package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type Static struct {
	Declaration *ASTNode
}

func (s *Static) Parse(lex []tokenizer.Token, i *int) (e error) {
	return parseAccessSpec(s, lex, i)
}

func (s *Static) SetDecl(node *ASTNode) {
	s.Declaration = node
}
