package obj

import (
	"fmt"

	"github.com/sushil-cmd-r/glox/vm/code"
)

type Function struct {
	name  string
	arity int

	code      []byte
	constants []Obj
}

func NewFunction(name string, arity int, code []byte, constants []Obj) *Function {
	fn := &Function{
		name:  name,
		arity: arity,

		code:      code,
		constants: constants,
	}

	return fn
}

func (f *Function) SetName(name string) {
	f.name = name
}

func (f *Function) Type() ObjType {
	return FuncObj
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.name)
}

func (f *Function) ReadInst(offset int) byte {
	return f.code[offset]
}

func (f *Function) ReadConstant(offset byte) Obj {
	return f.constants[offset]
}

func (f *Function) PrintCode() {
	name := f.name
	if name == "<init>" {
		name = "script"
	}
	fmt.Printf("<%s>\n", name)
	for i := 0; i < len(f.code); {
		i = f.printInst(i)
	}
}

func (f *Function) printInst(i int) int {
	b := f.code[i]
	inst := code.Opcodes[b]
	i += 1
	switch b {
	case code.OpConstant, code.OpGetGlobal, code.OpSetGlobal, code.OpDefineGlobal:
		idx := f.code[i]
		constant := f.constants[idx]
		i += 1
		fmt.Printf("%15s %s\n", inst, constant)

	case code.OpGetLocal, code.OpSetLocal, code.OpCall:
		idx := f.code[i]
		i += 1
		fmt.Printf("%15s %d\n", inst, idx)

	default:
		fmt.Printf("%15s\n", inst)
	}

	return i
}
