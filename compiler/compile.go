package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/initialize"
	"github.com/tusklang/tusk/varprocessor"
)

var processor = varprocessor.NewProcessor()

//list of all the variables that are added by default
//has types to begin with, but it can store anything
//types are variables in tusk's parser so we need to add the default ones in like so
var prevars = map[string]data.Value{
	"i128": data.NewPrimative(types.I128),
	"i64":  data.NewPrimative(types.I64),
	"i32":  data.NewPrimative(types.I32),
	"i16":  data.NewPrimative(types.I16),
	"i8":   data.NewPrimative(types.I8),
	"f64":  data.NewPrimative(types.Double),
	"f32":  data.NewPrimative(types.Float),
}

func Compile(prog *initialize.Program, outfile string) {

	var compiler ast.Compiler
	m := ir.NewModule() //create a new llvm module

	//set the module stuff in the compiler
	compiler.Module = m
	compiler.OS = runtime.GOOS
	compiler.ARCH = runtime.GOARCH

	var stringclass *data.Class
	stringclass, compiler.NewString = initString(m)
	prevars["string"] = stringclass

	compiler.VarMap = make(map[string]*data.Variable)

	initDefaultOps(&compiler)

	//initialize the operations
	var initfunc = m.NewFunc("_tusk_init", types.Void) //initialize func ran before main
	compiler.InitBlock = initfunc.NewBlock("")

	var (
		cpacks   = make(map[*initialize.Package]*data.Package)
		cclasses = make(map[*initialize.File]*data.Class)
	)

	for _, v := range prog.Packages {
		//create the new package
		tp := data.NewPackage(v.Name, v.FullName())
		cpacks[v] = tp //add it to the list of packages
	}

	//add all the classes (files) to the type list
	for _, v := range prog.Packages {

		packtyp := cpacks[v]
		parentPacktyp := cpacks[v.Parent()]

		if v.Parent() != nil {

			if v.Parent().Parent() == nil {
				//if it has no parent, it's a package on the uppermost level
				//so it is it's variable/type/thing

				//we check the parent's parent because the *real* uppermost level is the unnamed one
				prevars[v.Name] = packtyp
			}

			//also if it has a parent, we will append the child to the parent's children map
			parentPacktyp.ChildPacks[v.Name] = packtyp

		}

		for _, vv := range v.Files {

			tc := compileClass(&compiler, vv, v, packtyp)
			cclasses[vv] = tc

			if v.Parent() == nil {
				processor.AddPreDecl(vv.Name)
				compiler.AddVar(vv.Name, data.NewVariable(tc, tc, true))
			}
		}

	}

	//add all the prevars as variables to the compiler
	for k, v := range prevars {
		processor.AddPreDecl(k)
		compiler.AddVar(k, data.NewVariable(v, data.NewPrimative(v.Type()), true))
	}

	//process all the variables

	for _, v := range prog.Packages {
		for _, vv := range v.Files {
			processor.ProcessVars(vv)
		}
	}

	j, _ := json.MarshalIndent(prog, "", "  ")
	fmt.Println(string(j))

	for ic, c := range cclasses {
		for _, v := range ic.Globals {
			v.Value.DeclareGlobal(c.ParentPackage.FullName+"."+c.Name+"_"+v.Value.Name, &compiler, c, v.IsStatic)
		}
	}

	//init the fields in each class before compiling as well
	for ic, c := range cclasses {

		for _, v := range ic.Globals {

			if !v.IsStatic {
				//must be an instance to be included in the struct type
				c.SType.Fields = append(c.SType.Fields, v.Value.Type.Group.Compile(&compiler, c, v.Value.Type, c.Construct).Type())
			}

		}

	}

	for ic, c := range cclasses {

		var newAlloc = c.Construct.NewAlloca(c.SType)

		for _, v := range ic.Globals {

			if v.IsStatic {
				//it's a static variable
				v.Value.CompileGlobal(&compiler, c, compiler.InitBlock)
			} else {
				//it's an instance variable
				v.Value.CompileGlobal(&compiler, c, c.Construct)
			}

		}

		if ic.Constructor != nil {
			//add the constructor into the mix of this jazz
			ic.Constructor.CompileConstructor(&compiler, c, c.Construct, newAlloc)
		}

		loadedA := c.Construct.NewLoad(c.SType, newAlloc)
		c.Construct.NewRet(loadedA) //return void at the end of the instance constructor
	}

	//declare the llvm entry function
	mfunc := m.NewFunc("main", types.Void)
	mblock := mfunc.NewBlock("")

	compiler.InitBlock.NewRet(nil) //append a `return void` to the init function

	mblock.NewCall(initfunc) //call the initialize function
	// loaded := mblock.NewLoad(mfnc.ContentType, mfnc)
	// mblock.NewCall(loaded)
	mblock.NewRet(nil)

	f, _ := os.Create(outfile)
	f.Write([]byte(m.String()))
}
