package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	inst            Value
	Typ             Type
	UnReferenceable bool
}

func NewVariable(inst Value, typ Type, unReferenceable bool) *Variable {
	return &Variable{
		inst:            inst,
		Typ:             typ,
		UnReferenceable: unReferenceable,
	}
}

func (v *Variable) FetchAssig() Value {
	return v.inst
}

func (v *Variable) LLVal(block *ir.Block) value.Value {
	return block.NewLoad(v.Type(), v.inst.LLVal(block))
}

func (v *Variable) Type() types.Type {
	return v.Typ.Type()
}

func (v *Variable) TypeString() string {
	return v.Typ.TypeString()
}
