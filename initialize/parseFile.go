package initialize

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/tokenizer"
)

func parseFile(name string) (*File, error) {

	f, e := ioutil.ReadFile(name)

	if e != nil {
		return nil, e
	}

	lex := tokenizer.Tokenizer(string(f))
	generatedAST, e := ast.GenerateAST(lex)

	if e != nil {
		return nil, e
	}

	retf := File{
		Name: strings.TrimSuffix(filepath.Base(name), ".tusk"), //get the classname of the file
	}

	fetchGlobals(generatedAST, &retf, 2, false)

	//process variables
	varprocessor := NewProcessor()
	varprocessor.ProcessVars(&retf)
	///////////////////

	return &retf, nil
}
