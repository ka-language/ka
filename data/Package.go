package data

type Package struct {
	PackageName, FullName string
	Classes               map[string]*Class
	ChildPacks            map[string]*Package
}

func NewPackage(name, fullname string) *Package {
	return &Package{
		PackageName: name,
		FullName:    fullname,
		Classes:     make(map[string]*Class),
		ChildPacks:  make(map[string]*Package),
	}
}

func (p *Package) AddClass(name string, class *Class) {
	p.Classes[name] = class
}
