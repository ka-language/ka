package initialize

/*
	This interface represents either
		- package (folder)
		- class (file)
*/

type FSObj interface {
	Parent() FSObj
}
