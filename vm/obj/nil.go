package obj

type null struct{}

func Nil() *null {
	return &null{}
}

func (*null) Type() ObjType {
	return NilObj
}

func (*null) String() string {
	return "<nil>"
}
