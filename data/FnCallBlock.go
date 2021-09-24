package data

type FnCallBlock struct {
	Args []Value
}

func NewFnCallBlock() *FnCallBlock {
	return &FnCallBlock{}
}
