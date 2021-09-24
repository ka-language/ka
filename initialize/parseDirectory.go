package initialize

import (
	"io/ioutil"
	"path"
	"path/filepath"
)

func parseDirectory(dir string, pkg *Directory) error {

	fsinfo, e := ioutil.ReadDir(dir)

	if e != nil {
		return e
	}

	for _, v := range fsinfo {

		//joined path (of the parent directories and current one)
		jpth := path.Join(dir, v.Name())

		//current filesystem object (a fs obj is just a directory or a file)
		var cur FSObj

		if v.IsDir() {
			//a new package

			var spkg Directory

			//because the name is an array (see `Package.go`) we want to get the package names of all the parents
			spkg.Name = v.Name()

			spkg.parent = pkg //set the parent package

			e = parseDirectory(jpth, &spkg)

			if e != nil {
				return e
			}

			cur = &spkg

		} else {
			//a new class

			//only append a new class if it's a tusk file
			if filepath.Ext(v.Name()) != ".tusk" {
				continue
			}

			//a new class in the package
			pf, e := parseFile(jpth)

			//set the parent package of the file
			pf.pkg = pkg

			if e != nil {
				return e
			}

			cur = pf
		}

		pkg.Nested = append(pkg.Nested, cur)
	}

	return nil
}
