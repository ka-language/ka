package ast

import (
	"strings"

	"github.com/tusklang/tusk/tokenizer"
)

type String struct {
	dstring []byte
}

func (s *String) Parse(lex []tokenizer.Token, i *int) error {

	sv := lex[*i].Name

	sv = strings.TrimSuffix(strings.TrimPrefix(sv, "\""), "\"") //remove the leading and trailing quotes

	s.dstring = []byte(sv)

	return nil
}
