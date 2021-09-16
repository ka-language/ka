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
