package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Link struct {
	stname /*<- stored tname after varname mangling*/, TName, AName string
	DType                                                           *ASTNode
	Access                                                          int
}

func (l *Link) Parse(lex []tokenizer.Token, i *int) error {

	//format looks like
	//	link var tusk_name: fn() -> asm_name

	*i++

	if lex[*i].Name != "var" {
		//error
	}

	*i++

	if lex[*i].Type != "variable" {
		//must be the varname
	}

	tname := lex[*i].Name
	*i++

	if lex[*i].Name != ":" {
		//error
		//must supply type
	}

	*i++

	dtype, e := groupsToAST(groupSpecific(lex, i, nil, 1))

	if e != nil {
		return e
	}

	if lex[*i].Name != "->" {
		//error
	}

	*i++

	aname := lex[*i].Name

	l.TName = tname
	l.stname = tname
	l.AName = aname
	l.DType = dtype[0]
	l.Access = 2 //access is private by default

	return nil
}

func (l *Link) addToClass(lf *ir.Func, compiler *Compiler, dtype data.Value, class *data.Class) data.Value {
	tfd := data.NewLinkedFunc(lf, dtype.(*data.Function).RetType())
	tfd.SetLName(l.stname)
	compiler.AddVar(l.TName, tfd)

	class.AppendStatic(l.stname, tfd, tfd.TType(), l.Access)
	return nil
}

func (l *Link) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	aname := l.AName //name in the linked binary
	dtype := l.DType.Group.Compile(compiler, class, l.DType, function)

	if dtype.TypeData().Name() != "func" {
		//error
		//linked values must be functions
	}

	if lf, exists := compiler.LinkedFunctions[aname]; exists {
		return l.addToClass(lf, compiler, dtype, class)
	}

	dfunc := dtype.(*data.Function).LLFunc

	dfunc.SetName(aname)
	dfunc.Params = nil
	dfunc.Sig.Variadic = true        //make it a variadic function in case it is declared elsewhere
	dfunc.Sig.RetType = types.I64Ptr //make the return an i64 pointer, this way we can return any value, and cast it appropriately when called. the appropriate cast is provided by dtype.(*data.Function).RetType()

	compiler.LinkedFunctions[l.AName] = dfunc

	return l.addToClass(dfunc, compiler, dtype, class)
}
