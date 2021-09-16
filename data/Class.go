package data

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Class struct {
	Name      string
	SType     *types.StructType
	Instance  *linkedhashmap.Map
	Static    map[string]*Variable
	Construct *ir.Block

	ParentPackage *Package
}

func NewClass(name string, st *types.StructType, parent *Package) *Class {
	return &Class{
		Name:          name,
		SType:         st,
		ParentPackage: parent,
	}
}

func (c *Class) LLVal(block *ir.Block) value.Value {
	return nil
}

func (c *Class) Default() constant.Constant {
	return constant.NewNull(types.NewPointer(c.SType))
}

func (c *Class) Type() types.Type {
	return c.SType
}

func (c *Class) TypeString() string {
	return "class " + c.SType.Name()
}
