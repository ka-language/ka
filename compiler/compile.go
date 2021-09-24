package compiler

import (
	"github.com/tusklang/tusk/initialize"
	"github.com/tusklang/tusk/varprocessor"
)

//the processor used to predict variable types, validate types, and validate variable usages
var processor = varprocessor.NewProcessor()

func Compile(prog *initialize.Program, outfile string) {

}
