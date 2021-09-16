package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type Type interface {
	Type() types.Type
	TypeString() string
	Default() constant.Constant
}

func NewType(t types.Type) Type {
	switch t.(type) {
	case *types.FuncType:
		fin := Function{}
		fin.FTyp = t
		return &fin
	}
	return NewPrimative(t) //otherwise it's a primative
}
