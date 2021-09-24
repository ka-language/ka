package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type Group interface {
	Parse([]tokenizer.Token, *int) error
}
