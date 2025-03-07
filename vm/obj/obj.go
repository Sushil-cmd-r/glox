package obj

type Obj interface {
	Type() ObjType
	String() string
}

func New[V Obj, T any](val T, constructor func(T) V) V {
	return constructor(val)
}

func As[V Obj, T string | float64](v V, val func(V) T) T {
	return val(v)
}

type ObjType int

const (
	NumberObj ObjType = iota
	StringObj
	NilObj
	BoolObj
	FuncObj
)

var objTypes = [...]string{
	NumberObj: "number",
	StringObj: "string",
	NilObj:    "<nil>",
	BoolObj:   "boolean",
	FuncObj:   "FuncObj",
}

func (ot ObjType) String() string {
	return objTypes[ot]
}
