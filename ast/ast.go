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

type IdentExpr struct {
	Token token.Token
	Value string
}

// Return Stmt
type ReturnStmt struct {
	Token     token.Token
	ReturnVal Expr
}

func (ls *LetStmt) stmtNode() {}
func (ls *LetStmt) TokenLiteral() string {
	return ls.Token.Literal
}

func (i *IdentExpr) exprNode() {}
func (i *IdentExpr) TokenLiteral() string {
	return i.Token.Literal
}

func (rs *ReturnStmt) stmtNode() {}
func (rs *ReturnStmt) TokenLiteral() string {
	return rs.Token.Literal
}
