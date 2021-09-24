package initialize

import (
	"strconv"
)

/*
pub stat var main = fn() {
	var a = 43;
	{
		var b = 2;
	};
	{
		var b = 3;
	};
};

in the above example, there are two variables named `b` in different scopes

it would become

pub stat var main = fn() {
	var vd_1 = 43;
	{
		var vd_2 = 2;
	};
	{
		var vd_3 = 3;
	};
};

all of the variables' names got mangled, so there are no duplicates throughout the program
global variables are the only exception- because they exist throughout the whole file, so there will be no duplicated names of globals in diffrent scopes
*/

type VarProcessor struct {
	curvar  int
	predecl map[string]decl
}

func NewProcessor() VarProcessor {
	return VarProcessor{}
}

func (p *VarProcessor) nextvar() string {
	p.curvar++
	return "vd_" + strconv.Itoa(p.curvar)
}

func (p *VarProcessor) AddPreDecl(n string) {

	if p.predecl == nil {
		p.predecl = make(map[string]decl)
	}

	p.predecl[n] = decl{
		nname: n,
	}
}
