package ast

import (
	"fmt"
	"strings"

	"github.com/sushil-cmd-r/glox/token"
)

type Stmt interface {
	stmtNode()
}

type ExprStmt struct {
	Expression Expr
}

type LetStmt struct {
	Name  *IdentExpr
	Value Expr
}

type AssignStmt struct {
	Name  Expr
	Value Expr
}

type BlockStmt struct {
	Stmts []Stmt
}

type PrintStmt struct {
	Expr Expr
}

type FuncStmt struct {
	Name     *IdentExpr
	FuncExpr *FuncExpr
}

func (*ExprStmt) stmtNode()   {}
func (*LetStmt) stmtNode()    {}
func (*AssignStmt) stmtNode() {}
func (*BlockStmt) stmtNode()  {}
func (*PrintStmt) stmtNode()  {}
func (*FuncStmt) stmtNode()   {}

func (e *ExprStmt) String() string {
	return fmt.Sprintf("%s;\n", e.Expression)
}

func (l *LetStmt) String() string {
	return fmt.Sprintf("let %s = %s;\n", l.Name, l.Value)
}

func (a *AssignStmt) String() string {
	return fmt.Sprintf("%s = %s;\n", a.Name, a.Value)
}

func (b *BlockStmt) String() string {
	return b.toString(0)
}

func (b *BlockStmt) toString(depth int) string {
	var sb strings.Builder
	ident := strings.Repeat("  ", depth)
	sb.WriteString("{\n")

	for _, stmt := range b.Stmts {
		var str string

		switch stmt := stmt.(type) {
		case *BlockStmt:
			str = stmt.toString(depth + 1)
		default:
			str = fmt.Sprintf("%s", stmt)
		}

		sb.WriteString(ident + "  " + str)
	}
	sb.WriteString(ident + "}\n")

	return sb.String()
}

func (p *PrintStmt) String() string {
	return fmt.Sprintf("print %s;\n", p.Expr)
}

func (f *FuncStmt) String() string {
	return fmt.Sprintf("function %s%s\n", f.Name, f.FuncExpr.String()[3:])
}

type Expr interface {
	exprNode()
}

type BinaryExpr struct {
	Op    token.Token
	Left  Expr
	Right Expr
}

type UnaryExpr struct {
	Op   token.Token
	Left Expr
}

type GroupExpr struct {
	Expression Expr
}

type NumberLit struct {
	Value float64
}

type StringLit struct {
	Value string
}

type IdentExpr struct {
	Name string
}

type NilExpr struct{}

type CallExpr struct {
	Callee Expr
	Args   []Expr
}

type FuncExpr struct {
	Params []*IdentExpr
	Body   *BlockStmt
}

func (b *BinaryExpr) exprNode() {}
func (u *UnaryExpr) exprNode()  {}
func (g *GroupExpr) exprNode()  {}
func (n *NumberLit) exprNode()  {}
func (s *StringLit) exprNode()  {}
func (i *IdentExpr) exprNode()  {}
func (n *NilExpr) exprNode()    {}
func (c *CallExpr) exprNode()   {}
func (f *FuncExpr) exprNode()   {}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Op, b.Left, b.Right)
}

func (u *UnaryExpr) String() string {
	return fmt.Sprintf("(%s %s)", u.Op, u.Left)
}

func (g *GroupExpr) String() string {
	return fmt.Sprintf("%s", g.Expression)
}

func (n *NumberLit) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func (s *StringLit) String() string {
	return s.Value
}

func (i *IdentExpr) String() string {
	return i.Name
}

func (n *NilExpr) String() string {
	return "<nil>"
}

func (c *CallExpr) String() string {
	var args []string
	for _, a := range c.Args {
		args = append(args, fmt.Sprintf("%s", a))
	}

	return fmt.Sprintf("%s(%s)", c.Callee, strings.Join(args, ","))
}

func (f *FuncExpr) String() string {
	var params []string
	for _, p := range f.Params {
		params = append(params, fmt.Sprintf("%s", p))
	}

	ps := strings.Join(params, ",")
	var sb strings.Builder

	sb.WriteString("fn (" + ps + ") ")
	sb.WriteString(f.Body.String())

	return sb.String()[:sb.Len()-1]
}
