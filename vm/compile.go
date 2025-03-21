package vm

import (
	"errors"
	"fmt"
	"math"

	"github.com/sushil-cmd-r/glox/ast"
	"github.com/sushil-cmd-r/glox/token"
	"github.com/sushil-cmd-r/glox/vm/code"
	"github.com/sushil-cmd-r/glox/vm/obj"
)

type Compiler struct {
	code      []byte
	constants []obj.Obj

	locals [math.MaxUint8]Local

	localCount int
	scopeDepth int

	debug bool
}

func initCompiler() *Compiler {
	c := &Compiler{
		code:      make([]byte, 0),
		constants: make([]obj.Obj, 0),

		localCount: 0,
		scopeDepth: 0,
	}

	l := Local{depth: 0, name: ""}
	c.locals[c.localCount] = l
	c.localCount += 1

	return c
}

var annonymousCnt = 0

func Compile(prog *ast.FuncExpr, fname string, debug bool) (*obj.Function, error) {
	c := initCompiler()
	c.debug = debug

	if fname != InitFunc {
		c.beginScope()
	}

	for _, p := range prog.Params {
		arg, err := c.registerDeclaration(p)
		if err != nil {
			return nil, err
		}

		if err := c.defineVariable(arg); err != nil {
			return nil, err
		}
	}

	for _, stmt := range prog.Body.Stmts {
		if err := c.compileStmt(stmt); err != nil {
			return nil, err
		}
	}

	c.emitInst(code.OpReturn, nil)
	if fname != InitFunc {
		c.endScope()
	}

	fn := obj.NewFunction(fname, len(prog.Params), c.code, c.constants)

	if debug {
		fn.PrintCode()
	}
	return fn, nil
}

func (c *Compiler) compileStmt(stmt ast.Stmt) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return c.compileExprStmt(stmt)

	case *ast.LetStmt:
		return c.compileLetStmt(stmt)

	case *ast.AssignStmt:
		return c.compileAssignStmt(stmt)

	case *ast.BlockStmt:
		return c.compileBlockStmt(stmt)

	case *ast.PrintStmt:
		return c.compilePrintStmt(stmt)

	case *ast.FuncStmt:
		return c.compileFuncStmt(stmt)

	default:
		panic("unimplemented stmt")
	}
}

func (c *Compiler) compileFuncStmt(stmt *ast.FuncStmt) error {
	arg, err := c.registerDeclaration(stmt.Name)
	if err != nil {
		return err
	}

	ident := stmt.Name
	fn, err := Compile(stmt.FuncExpr, ident.Name, c.debug)
	if err != nil {
		return err
	}

	c.emitInst(code.OpConstant, fn)
	c.defineVariable(arg)
	return nil
}

func (c *Compiler) compilePrintStmt(stmt *ast.PrintStmt) error {
	if err := c.compileExpr(stmt.Expr); err != nil {
		return err
	}

	c.emitInst(code.OpPrint, nil)
	return nil
}

func (c *Compiler) compileBlockStmt(stmt *ast.BlockStmt) error {
	c.beginScope()

	for _, stmt := range stmt.Stmts {
		if err := c.compileStmt(stmt); err != nil {
			return err
		}
	}

	c.endScope()
	return nil
}

func (c *Compiler) beginScope() {
	c.scopeDepth += 1
}

func (c *Compiler) endScope() {
	c.scopeDepth -= 1

	for i := c.localCount - 1; i >= 0 && c.locals[i].depth > c.scopeDepth; i-- {
		c.emitInst(code.OpPop, nil)
		c.localCount -= 1
	}
}

func (c *Compiler) compileAssignStmt(stmt *ast.AssignStmt) error {
	if err := c.compileExpr(stmt.Value); err != nil {
		return err
	}

	ident, ok := stmt.Name.(*ast.IdentExpr)
	if !ok {
		return fmt.Errorf("invalid assignment target: %T", stmt.Name)
	}

	arg := c.resolveLocal(ident.Name)
	if arg == -1 {
		name := obj.NewStr(ident.Name)
		i, err := c.addConstant(name)
		if err != nil {
			return err
		}
		c.emitInsts(code.OpSetGlobal, i)
	} else {
		c.emitInsts(code.OpSetLocal, byte(arg))
	}

	return nil
}

func (c *Compiler) compileLetStmt(stmt *ast.LetStmt) error {
	i, err := c.registerDeclaration(stmt.Name)
	if err != nil {
		return err
	}

	if err := c.compileExpr(stmt.Value); err != nil {
		return err
	}

	return c.defineVariable(i)
}

func (c *Compiler) registerDeclaration(ident *ast.IdentExpr) (byte, error) {
	if err := c.declareVariable(ident.Name); err != nil {
		return 0, err
	}

	if c.scopeDepth > 0 {
		return 0, nil
	}

	str := obj.NewStr(ident.Name)
	return c.addConstant(str)
}

type Local struct {
	name  string
	depth int
}

func (c *Compiler) declareVariable(name string) error {
	if c.scopeDepth == 0 {
		return nil
	}

	for i := c.localCount - 1; i >= 0; i-- {
		l := c.locals[i]
		if l.depth != -1 && l.depth < c.scopeDepth {
			break
		}

		if name == l.name {
			return fmt.Errorf("redeclaration of variable: %s", name)
		}
	}

	l := Local{name: name, depth: c.scopeDepth}
	c.locals[c.localCount] = l
	c.localCount += 1
	return nil
}

func (c *Compiler) defineVariable(i byte) error {
	if c.scopeDepth > 0 {
		return nil
	}

	c.emitInsts(code.OpDefineGlobal, i)
	return nil
}

func (c *Compiler) compileExprStmt(stmt *ast.ExprStmt) error {
	if err := c.compileExpr(stmt.Expression); err != nil {
		return err
	}

	c.emitInst(code.OpPop, nil)
	return nil
}

func (c *Compiler) compileExpr(expr ast.Expr) error {
	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		return c.compileBinary(expr)

	case *ast.UnaryExpr:
		return c.compileUnary(expr)

	default:
		return c.compilePrimary(expr)
	}
}

func (c *Compiler) compileBinary(expr *ast.BinaryExpr) error {
	if err := c.compileExpr(expr.Left); err != nil {
		return err
	}

	if err := c.compileExpr(expr.Right); err != nil {
		return err
	}

	var err error

	switch expr.Op {
	case token.PLUS:
		err = c.emitInst(code.OpAdd, nil)
	case token.MINUS:
		err = c.emitInst(code.OpSub, nil)
	case token.STAR:
		err = c.emitInst(code.OpMul, nil)
	case token.SLASH:
		err = c.emitInst(code.OpDiv, nil)
	}

	return err
}

func (c *Compiler) compileUnary(expr *ast.UnaryExpr) error {
	if err := c.compileExpr(expr.Left); err != nil {
		return err
	}

	var err error
	switch expr.Op {
	case token.MINUS:
		err = c.emitInst(code.OpNegate, nil)
	case token.NOT:
		err = c.emitInst(code.OpNot, nil)
	}
	return err
}

func (c *Compiler) compilePrimary(expr ast.Expr) error {
	switch expr := expr.(type) {
	case *ast.NumberLit:
		num := obj.NewNumber(expr.Value)
		return c.emitInst(code.OpConstant, num)

	case *ast.GroupExpr:
		return c.compileExpr(expr.Expression)

	case *ast.NilExpr:
		return c.emitInst(code.OpNil, nil)

	case *ast.IdentExpr:
		return c.compileIdent(expr)

	case *ast.CallExpr:
		return c.compileCallExpr(expr)

	case *ast.FuncExpr:
		fn, err := Compile(expr, fmt.Sprintf("Annonymous:%d", annonymousCnt), c.debug)
		annonymousCnt += 1
		if err != nil {
			return err
		}
		return c.emitInst(code.OpConstant, fn)

	default:
		fmt.Printf("%T\n", expr)
		panic("not implemented yet")
	}
}

func (c *Compiler) compileCallExpr(expr *ast.CallExpr) error {
	if err := c.compileExpr(expr.Callee); err != nil {
		return err
	}

	for _, a := range expr.Args {
		if err := c.compileExpr(a); err != nil {
			return err
		}
	}

	c.emitInsts(code.OpCall, byte(len(expr.Args)))
	return nil
}

func (c *Compiler) compileIdent(expr *ast.IdentExpr) error {
	arg := c.resolveLocal(expr.Name)
	if arg == -1 {
		str := obj.NewStr(expr.Name)
		i, err := c.addConstant(str)
		if err != nil {
			return err
		}
		c.emitInsts(code.OpGetGlobal, i)

	} else {
		c.emitInsts(code.OpGetLocal, byte(arg))
	}

	return nil
}

func (c *Compiler) resolveLocal(name string) int {
	for i := c.localCount - 1; i >= 0; i-- {
		if c.locals[i].name == name {
			return i
		}
	}

	return -1
}

func (c *Compiler) emitInst(opcode byte, o obj.Obj) error {
	c.code = append(c.code, opcode)

	if o != nil {
		idx, err := c.addConstant(o)
		if err != nil {
			return err
		}
		c.code = append(c.code, idx)
	}

	return nil
}

func (c *Compiler) emitInsts(o1, o2 byte) {
	c.code = append(c.code, o1, o2)
}

var ErrTooManyconstants = errors.New("too many constants")

func (c *Compiler) addConstant(o obj.Obj) (byte, error) {
	if len(c.constants) == math.MaxUint8 {
		return 0, ErrTooManyconstants
	}
	c.constants = append(c.constants, o)
	return byte(len(c.constants) - 1), nil
}
