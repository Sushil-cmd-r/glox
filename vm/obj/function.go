package obj

import (
	"errors"
	"fmt"
	"math"

	"github.com/sushil-cmd-r/glox/vm/code"
)

type FuncType int

const (
	Script FuncType = iota
	Funct
)

type Func struct {
	fType FuncType
	name  string

	arity int

	code      []byte
	constants []Obj
}

func Function(fType FuncType) *Func {
	fn := &Func{
		fType: fType,

		code:      make([]byte, 0),
		constants: make([]Obj, 0, math.MaxUint8),
	}

	return fn
}

func (f *Func) SetName(name string) {
	f.name = name
}

func (f *Func) Type() ObjType {
	return FuncObj
}

func (f *Func) String() string {
	return fmt.Sprintf("<fn %s>", f.name)
}

func (f *Func) ReadInst(offset int) byte {
	return f.code[offset]
}

func (f *Func) ReadConstant(offset byte) Obj {
	return f.constants[offset]
}

func (f *Func) EmitInst(opcode byte, o Obj) error {
	f.code = append(f.code, opcode)

	if o != nil {
		idx, err := f.AddConstant(o)
		if err != nil {
			return err
		}
		f.code = append(f.code, idx)
	}

	return nil
}

func (f *Func) EmitInsts(o1, o2 byte) {
	f.code = append(f.code, o1, o2)
}

var ErrTooManyconstants = errors.New("too many constants")

func (f *Func) AddConstant(o Obj) (byte, error) {
	if len(f.constants) == math.MaxUint8 {
		return 0, ErrTooManyconstants
	}
	f.constants = append(f.constants, o)
	return byte(len(f.constants) - 1), nil
}

func (f *Func) PrintCode() {
	name := f.name
	if name == "" {
		name = "script"
	}
	fmt.Printf("<%s>\n", name)
	for i := 0; i < len(f.code); {
		i = f.printInst(i)
	}
}

func (f *Func) printInst(i int) int {
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
