package compiler

import (
	"fmt"
	"os"

	. "github.com/omm-lang/omm/lang/types"

	. "github.com/omm-lang/omm/lang/interpreter"
)

var included = []string{} //list of the imported files from omm

func Run(params CliParams) {

	fileName := params.Name

	included = append(included, fileName)

	variables, ce := Compile(params)

	if ce != nil {
		fmt.Println(ce.Error())
		os.Exit(1)
	}

	RunInterpreter(variables, params)
}
