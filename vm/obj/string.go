package obj

type Str struct {
	value string
}

func NewStr(val string) *Str {
	return &Str{value: val}
}

func (s *Str) Type() ObjType {
	return StringObj
}

func AsStr(o Obj) string {
	return o.(*Str).value
}

func (s *Str) String() string {
	return s.value
}
