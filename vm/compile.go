package vm

import (
	"fmt"
	"math"

	"github.com/sushil-cmd-r/glox/ast"
	"github.com/sushil-cmd-r/glox/token"
	"github.com/sushil-cmd-r/glox/vm/obj"
)

type Compiler struct {
	function *obj.Func

	locals [math.MaxUint8]Local

	localCount int
	scopeDepth int
}

func initCompiler(t obj.FuncType) *Compiler {
	c := &Compiler{
		function: obj.New(t, obj.Function),
	}

	l := Local{depth: 0, name: ""}
	c.locals[c.localCount] = l
	c.localCount += 1

	return c
}

func (c *Compiler) Compile(prog []ast.Stmt) (*obj.Func, error) {
	if err := c.compile(prog); err != nil {
		return nil, fmt.Errorf("Compilation Error: %w", err)
	}

	return c.function, nil
}

func (c *Compiler) compile(prog []ast.Stmt) error {
	for _, stmt := range prog {
		if err := c.compileStmt(stmt); err != nil {
			return err
		}
	}
	c.emitInst(OpReturn, nil)
	return nil
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

	c1 := initCompiler(obj.Funct)
	c1.function.SetName(stmt.Name.String())
	c1.beginScope()
	for _, p := range stmt.Params {
		arg, err := c1.registerDeclaration(p)
		if err != nil {
			return err
		}

		c1.defineVariable(arg)
	}

	c1.compileBlockStmt(stmt.Body)
	c1.emitInst(OpReturn, nil)
	c1.endScope()

	c.emitInst(OpConstant, c1.function)

	c.defineVariable(arg)
	return nil
}

func (c *Compiler) compilePrintStmt(stmt *ast.PrintStmt) error {
	if err := c.compileExpr(stmt.Expr); err != nil {
		return err
	}

	c.emitInst(OpPrint, nil)
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
		c.function.EmitInst(OpPop, nil)
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
		name := obj.New(ident.Name, obj.Str)
		i, err := c.addConstant(name)
		if err != nil {
			return err
		}
		c.emitInsts(OpSetGlobal, i)
	} else {
		c.emitInsts(OpSetLocal, byte(arg))
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

	str := obj.New(ident.Name, obj.Str)
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

	c.emitInsts(OpDefineGlobal, i)
	return nil
}

func (c *Compiler) compileExprStmt(stmt *ast.ExprStmt) error {
	if err := c.compileExpr(stmt.Expression); err != nil {
		return err
	}

	c.emitInst(OpPop, nil)
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
		err = c.emitInst(OpAdd, nil)
	case token.MINUS:
		err = c.emitInst(OpSub, nil)
	case token.STAR:
		err = c.emitInst(OpMul, nil)
	case token.SLASH:
		err = c.emitInst(OpDiv, nil)
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
		err = c.emitInst(OpNegate, nil)
	case token.NOT:
		err = c.emitInst(OpNot, nil)
	}
	return err
}

func (c *Compiler) compilePrimary(expr ast.Expr) error {
	switch expr := expr.(type) {
	case *ast.NumberLit:
		num := obj.New(expr.Value, obj.Number)
		return c.emitInst(OpConstant, num)

	case *ast.GroupExpr:
		return c.compileExpr(expr.Expression)

	case *ast.NilExpr:
		return c.emitInst(OpNil, nil)

	case *ast.IdentExpr:
		return c.compileIdent(expr)

	case *ast.CallExpr:
		return c.compileCallExpr(expr)

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

	c.emitInsts(OpCall, byte(len(expr.Args)))
	return nil
}

func (c *Compiler) compileIdent(expr *ast.IdentExpr) error {
	arg := c.resolveLocal(expr.Name)
	if arg == -1 {
		str := obj.New(expr.Name, obj.Str)
		i, err := c.addConstant(str)
		if err != nil {
			return err
		}
		c.emitInsts(OpGetGlobal, i)

	} else {
		c.emitInsts(OpGetLocal, byte(arg))
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
	return c.function.EmitInst(opcode, o)
}

func (c *Compiler) emitInsts(o1, o2 byte) {
	c.function.EmitInsts(o1, o2)
}

func (c *Compiler) addConstant(o obj.Obj) (byte, error) {
	return c.function.AddConstant(o)
}
