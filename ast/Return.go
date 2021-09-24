package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type Return struct {
	Val *ASTNode
}

func (r *Return) Parse(lex []tokenizer.Token, i *int) error {

	*i++

	retval := braceMatcher(lex, i, allopeners, allclosers, false, "terminator")
	retvalAST, e := groupsToAST(grouper(retval))

	r.Val = retvalAST[0]

	return e
}
