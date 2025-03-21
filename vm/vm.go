package vm

import (
	"fmt"
	"math"

	"github.com/sushil-cmd-r/glox/ast"
	"github.com/sushil-cmd-r/glox/parser"
	"github.com/sushil-cmd-r/glox/vm/obj"
)

const InitFunc = "<init>"

type CalLFrame struct {
	function *obj.Function
	ip       int
	stack    []obj.Obj
}

type VM struct {
	fp        int
	frames    [64]*CalLFrame
	currFrame *CalLFrame

	sp    int
	stack [64 * math.MaxUint8]obj.Obj

	globals map[string]obj.Obj
	debug   bool
}

func Init(debug bool) *VM {
	vm := &VM{fp: -1, globals: make(map[string]obj.Obj), debug: debug}
	return vm
}

func (vm *VM) Execute(src []byte) error {
	p := parser.New(src)
	prog, err := p.Parse()
	if err != nil {
		return fmt.Errorf("Syntax Error: %w", err)
	}

	fnExpr := &ast.FuncExpr{Params: nil, Body: &ast.BlockStmt{Stmts: prog}}
	function, err := Compile(fnExpr, InitFunc, vm.debug)
	if err != nil {
		return fmt.Errorf("Compilation Error: %w", err)
	}

	vm.push(function)
	vm.call(function, 0)
	// return nil
	return vm.run()
}

func (vm *VM) call(o obj.Obj, args byte) error {
	if o.Type() != obj.FuncObj {
		return fmt.Errorf("%s not callable", o.Type())
	}

	fn := o.(*obj.Function)
	frame := &CalLFrame{
		ip:       0,
		function: fn,
		stack:    vm.stack[vm.sp-int(args)-1:],
	}

	vm.fp += 1
	vm.currFrame = frame
	vm.frames[vm.fp] = frame

	return nil
}
