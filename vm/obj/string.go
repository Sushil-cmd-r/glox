package obj

import "fmt"

type str struct {
	value string
}

func Str(val string) *str {
	return &str{value: val}
}

func (s *str) Type() ObjType {
	return StringObj
}

func StrVal(o Obj) string {
	return o.(*str).value
}

func (s *str) String() string {
	return fmt.Sprintf("String{value: %v}", s.value)
}
