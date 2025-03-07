package obj

import "fmt"

type number struct {
	value float64
}

func Number(val float64) *number {
	return &number{value: val}
}

func (n *number) Type() ObjType {
	return NumberObj
}

func NumVal(o Obj) float64 {
	return o.(*number).value
}

func (n *number) String() string {
	return fmt.Sprintf("Number{value: %v}", n.value)
}
