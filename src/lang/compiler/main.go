package compiler

import "io/ioutil"
import "fmt"
import "os"
import "strings"

import . "lang/types"
import . "lang/interpreter"

var included = []string{} //list of the imported files from omm
var ommbasedir string //directory of the omm installation

//export Run
func Run(params CliParams) {

  ommbasedir = params.OmmDirname
  fileName := params.Name

  var compileall = false
  if strings.HasSuffix(fileName, "*") || strings.HasSuffix(fileName, "*/") {
    compileall = true
    fileName = "main.omm"
  }

  included = append(included, fileName)

  file, e := ioutil.ReadFile(fileName)

  if e != nil {
    fmt.Println("Could not find", fileName)
    os.Exit(1)
  }

  _, variables, ce := Compile(string(file), fileName, compileall, true)

  if ce != nil {
    ce.Print()
  }

  RunInterpreter(variables, params)
}
