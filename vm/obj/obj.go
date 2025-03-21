package obj

type Obj interface {
	Type() ObjType
	String() string
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
