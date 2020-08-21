package types

import "strconv"

type OmmBool struct {
	Boolean *bool
}

func (b *OmmBool) FromGoType(val bool) {
	b.Boolean = &val
}

func (b OmmBool) ToGoType() bool {

	if b.Boolean == nil {
		return false
	}

	return *b.Boolean
}

func (b OmmBool) Format() string {
	return strconv.FormatBool(*b.Boolean)
}

func (b OmmBool) Type() string {
	return "bool"
}

func (b OmmBool) TypeOf() string {
	return b.Type()
}

func (_ OmmBool) Deallocate() {}