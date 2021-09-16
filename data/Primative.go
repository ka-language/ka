package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Primative struct {
	typ  types.Type
	name *string
}

func NewPrimative(typ types.Type) *Primative {
	return &Primative{
		typ: typ,
	}
}

func (p *Primative) Default() constant.Constant {
	switch p.typ {
	case types.I32:
		return constant.NewInt(types.I32, 0)
	default:
		return &constant.Null{}
	}
}

func (p *Primative) SetTypeName(n string) {
	p.name = &n
}

func (p *Primative) LLVal(block *ir.Block) value.Value {
	return nil
}

func (p *Primative) Type() types.Type {
	return p.typ
}

func (p *Primative) typMatch() string {

	if p.name != nil {
		return *p.name
	}

	return p.Type().LLString()
}

func (p *Primative) TypeString() string {
	return p.typMatch()
}
