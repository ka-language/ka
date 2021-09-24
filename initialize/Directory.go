package initialize

type Directory struct {
	Name   string
	Nested []FSObj
	parent *Directory
}

func (p *Directory) Parent() FSObj {
	return p.parent
}

func (p Directory) FullName() string {
	//returns the full name, with all the parents
	if p.parent == nil {
		return p.Name
	}
	return p.parent.FullName() + p.Name + "."
}
