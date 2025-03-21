package obj

import "fmt"

type Number struct {
	value float64
}

func NewNumber(val float64) *Number {
	return &Number{value: val}
}

func (n *Number) Type() ObjType {
	return NumberObj
}

func AsNum(o Obj) float64 {
	return o.(*Number).value
}

func (n *Number) String() string {
	return fmt.Sprintf("%v", n.value)
}
