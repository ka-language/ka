package data

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Class struct {
	Name      string
	Instance  *linkedhashmap.Map
	Static    *linkedhashmap.Map
	Construct *ir.Block

	ParentPackage *Package
}

func NewClass(name string, st *types.StructType, parent *Package) *Class {
	return &Class{
		Name:          name,
		ParentPackage: parent,
		Instance:      linkedhashmap.New(),
		Static:        linkedhashmap.New(),
	}
}

func (c *Class) LLVal(block *ir.Block) value.Value {
	return nil
}
