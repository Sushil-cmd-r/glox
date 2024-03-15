package ast

import (
	"github.com/sushil-cmd-r/glox/token"
)

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
	Stringify(depth int) []byte
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

// Expression Stmt
type ExpressionStmt struct {
	Token      token.Token
	Expression Expr
}

// Number Literal
type NumberLiteral struct {
	Token token.Token
	Value float64
}

// Prefix Expr
type PrefixExpr struct {
	Token    token.Token
	Operator string
	Right    Expr
}

// Infix Expr
type InfixExpr struct {
	Token    token.Token
	Left     Expr
	Operator string
	Right    Expr
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

func (es *ExpressionStmt) stmtNode() {}
func (es *ExpressionStmt) TokenLiteral() string {
	return es.Token.Literal
}

func (nl *NumberLiteral) exprNode() {}
func (nl *NumberLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (pe *PrefixExpr) exprNode() {}
func (pe *PrefixExpr) TokenLiteral() string {
	return pe.Token.Literal
}

func (ie *InfixExpr) exprNode() {}
func (ie *InfixExpr) TokenLiteral() string {
	return ie.Token.Literal
}
