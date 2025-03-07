package vm

import (
	"fmt"

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
		case OpReturn:
			return nil

		case OpConstant:
			obj := vm.readConstant()
			vm.push(obj)

		case OpAdd, OpSub, OpMul, OpDiv:
			err = vm.binaryOp(op)

		case OpNegate, OpNot:
			err = vm.unaryOp(op)

		case OpNil:
			vm.push(obj.Nil())

		case OpPop:
			vm.pop()

		case OpPrint:
			fmt.Println(vm.pop())

		case OpDefineGlobal:
			oj := vm.readConstant()
			name := obj.As(oj, obj.StrVal)
			vm.globals[name] = vm.pop()

		case OpGetGlobal:
			oj := vm.readConstant()
			name := obj.As(oj, obj.StrVal)
			obj, ok := vm.globals[name]
			if !ok {
				return fmt.Errorf("undefined variable: %s", name)
			}
			vm.push(obj)

		case OpSetGlobal:
			oj := vm.readConstant()
			name := obj.As(oj, obj.StrVal)
			vm.globals[name] = vm.pop()

		case OpGetLocal:
			i := vm.readInst()
			vm.push(vm.currFrame.stack[i])

		case OpSetLocal:
			i := vm.readInst()
			vm.currFrame.stack[i] = vm.pop()

		case OpCall:
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
	if op != OpAdd {
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
	case OpAdd:
		res = na + nb
	case OpSub:
		res = na - nb
	case OpMul:
		res = na * nb
	case OpDiv:
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
	return vm.currFrame.function.ReadInst()
}

func (vm *VM) readConstant() obj.Obj {
	idx := vm.readInst()
	return vm.currFrame.function.ReadConstant(idx)
}
