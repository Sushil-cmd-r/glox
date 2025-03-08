package vm

import (
	"fmt"

	"github.com/sushil-cmd-r/glox/vm/code"
	"github.com/sushil-cmd-r/glox/vm/obj"
)

// func (vm *VM) Interpret() error {
//
// 	return vm.run()
// }

func (vm *VM) run() error {
	for {
		op := vm.readInst()

		var err error
		switch op {
		case code.OpReturn:
			vm.fp -= 1
			if vm.fp == -1 {
				vm.pop()
				return nil
			}
			vm.currFrame = vm.frames[vm.fp]

		case code.OpConstant:
			obj := vm.readConstant()
			vm.push(obj)

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err = vm.binaryOp(op)

		case code.OpNegate, code.OpNot:
			err = vm.unaryOp(op)

		case code.OpNil:
			vm.push(obj.Nil())

		case code.OpPop:
			vm.pop()

		case code.OpPrint:
			fmt.Println(vm.pop())

		case code.OpDefineGlobal:
			oj := vm.readConstant()
			name := obj.As(oj, obj.StrVal)
			vm.globals[name] = vm.pop()

		case code.OpGetGlobal:
			oj := vm.readConstant()
			name := obj.As(oj, obj.StrVal)
			obj, ok := vm.globals[name]
			if !ok {
				return fmt.Errorf("undefined variable: %s", name)
			}
			vm.push(obj)

		case code.OpSetGlobal:
			oj := vm.readConstant()
			name := obj.As(oj, obj.StrVal)
			vm.globals[name] = vm.pop()

		case code.OpGetLocal:
			i := vm.readInst()
			vm.push(vm.currFrame.stack[i])

		case code.OpSetLocal:
			i := vm.readInst()
			vm.currFrame.stack[i] = vm.pop()

		case code.OpCall:
			args := vm.readInst()
			if err := vm.call(vm.stack[vm.sp-int(args)-1], args); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

func (vm *VM) unaryOp(_ byte) error {
	a := vm.pop()

	if a.Type() != obj.NumberObj {
		return fmt.Errorf("invalid unary operation on %s", a.Type())
	}

	val := obj.As(a, obj.NumVal)
	vm.push(obj.New(val*-1, obj.Number))

	return nil
}

func (vm *VM) binaryOp(op byte) error {
	b := vm.pop()
	a := vm.pop()

	if a.Type() != b.Type() {
		return fmt.Errorf("invalid operation between %s and %s", a.Type(), b.Type())
	}

	switch a.Type() {
	case obj.NumberObj:
		return vm.binaryOpNumber(op, a, b)

	case obj.StringObj:
		return vm.binaryOpString(op, a, b)

	default:
		return fmt.Errorf("invalid operation %d between %s", op, a.Type())
	}
}

func (vm *VM) binaryOpString(op byte, a, b obj.Obj) error {
	if op != code.OpAdd {
		return fmt.Errorf("invalid operation %d between strings", op)
	}

	sa := obj.As(a, obj.StrVal)
	sb := obj.As(b, obj.StrVal)

	obj := obj.New(sa+sb, obj.Str)
	vm.push(obj)

	return nil
}

func (vm *VM) binaryOpNumber(op byte, a, b obj.Obj) error {
	na := obj.As(a, obj.NumVal)
	nb := obj.As(b, obj.NumVal)

	var res float64
	switch op {
	case code.OpAdd:
		res = na + nb
	case code.OpSub:
		res = na - nb
	case code.OpMul:
		res = na * nb
	case code.OpDiv:
		if nb == 0 {
			return fmt.Errorf("division by zero")
		}
		res = na / nb
	}

	obj := obj.New(res, obj.Number)
	vm.push(obj)

	return nil
}

func (vm *VM) push(obj obj.Obj) {
	if vm.sp == len(vm.stack) {
		panic("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp += 1
}

func (vm *VM) pop() obj.Obj {
	if vm.sp <= 0 {
		panic("stack underflow")
	}
	vm.sp -= 1
	return vm.stack[vm.sp]
}

func (vm *VM) readInst() byte {
	offset := vm.currFrame.ip
	vm.currFrame.ip += 1
	return vm.currFrame.function.ReadInst(offset)
}

func (vm *VM) readConstant() obj.Obj {
	idx := vm.readInst()
	return vm.currFrame.function.ReadConstant(idx)
}
