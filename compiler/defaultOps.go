package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

func initDefaultOps(compiler *ast.Compiler) {

	compiler.OperationStore = ast.NewOperationStore()

	compiler.OperationStore.NewOperation("+", "i32", "i32", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		return data.NewInstruction((block.NewAdd(left.LLVal(block), right.LLVal(block))))
	})

	compiler.OperationStore.NewOperation(".", "package", "udvar", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		pack := left.(*data.Package)
		sub := right.(*data.UndeclaredVar).Name

		//it can either be a class or a subpackage
		var (
			class   = pack.Classes[sub]
			subpack = pack.ChildPacks[sub]
		)

		if class == nil {
			return subpack
		}

		return class
	})

	compiler.OperationStore.NewOperation(".", "class", "udvar", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		sub := right.(*data.UndeclaredVar).Name

		switch class := left.(type) {
		case *data.Class:
			//accessing static portion of the class

			return class.Static[sub]
		case *data.Variable:

			ctyp := class.Typ.(*data.Class)
			keys := ctyp.Instance.Keys()

			var idx int64 //store the index of the class' field

			for k, v := range keys {
				if v.(string) == sub {
					idx = int64(k)
					break
				}
			}

			//use GEP to fetch the field
			inst := block.NewGetElementPtr(
				ctyp.SType,
				class.FetchAssig().LLVal(block),
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, idx),
			)

			return data.NewInstruction(inst)
		}

		return nil
	})

	compiler.OperationStore.NewOperation("()", "func", "fncallb", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		f := left.LLVal(block)
		fcb := right.(*data.FnCallBlock)

		var args []value.Value

		for _, v := range fcb.Args {
			args = append(args, v.LLVal(block))
		}

		return data.NewInstruction(
			block.NewCall(f, args...),
		)
	})

	compiler.OperationStore.NewOperation("()", "class", "fncallb", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		class := left.(*data.Class)
		fcb := right.(*data.FnCallBlock)

		var args []value.Value

		for _, v := range fcb.Args {
			args = append(args, v.LLVal(block))
		}

		return data.NewInstruction(
			block.NewCall(class.Construct.Parent, args...),
		)
	})

	compiler.OperationStore.NewOperation("=", "ptr", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		ptr := left.LLVal(block)
		block.NewStore(right.LLVal(block), ptr)
		return nil
	})

}
