package ast

import "github.com/sushil-cmd-r/glox/token"

type Node interface {
	TokenLiteral() string
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLiteral()
	}
	return ""
}

// Let Stmt
type LetStmt struct {
	Token token.Token
	Name  *IdentExpr
	Value Expr
}

func (ls *LetStmt) stmtNode() {}
func (ls *LetStmt) TokenLiteral() string {
	return ls.Token.Literal
}

type IdentExpr struct {
	Token token.Token
	Value string
}

func (i *IdentExpr) exprNode() {}
func (i *IdentExpr) TokenLiteral() string {
	return i.Token.Literal
}
