package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Function struct {
	llfunc *ir.Func
	ftype  types.Type
}

func NewFunc(f *ir.Func) *Function {
	return &Function{
		llfunc: f,
		ftype:  f.Type(),
	}
}

func (f *Function) LLVal(block *ir.Block) value.Value {
	return f.llfunc
}

func (f *Function) Default() constant.Constant {
	return constant.NewNull(types.NewPointer(f.ftype))
}

func (f *Function) Type() types.Type {
	return f.llfunc.Type()
}

func (f *Function) TypeString() string {
	return "func"
}
