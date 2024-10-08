package parser

import (
	"fmt"
	"strconv"

	"github.com/sushil-cmd-r/glox/ast"
	"github.com/sushil-cmd-r/glox/scanner"
	"github.com/sushil-cmd-r/glox/token"
)

type Parser struct {
	sc     *scanner.Scanner
	file   *token.File
	errors *ErrorList

	tok token.Token
	lit string
	loc token.Loc
}

func New(file *token.File, source string) *Parser {
	sc := scanner.Init(file, []byte(source))
	p := &Parser{sc: sc, file: file, errors: &ErrorList{}}

	p.advance()
	return p
}

func (p *Parser) Parse() []ast.Stmt {
	stmts := make([]ast.Stmt, 0)
	for p.tok != token.EOF {
		stmt := p.parseStmt()
		p.expectSemi()
		stmts = append(stmts, stmt)
	}
	return stmts
}

func (p *Parser) parseStmt() ast.Stmt {
	return p.parseExprStmt()
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{}
	stmt.Expression = p.parseExpr(token.PrecLowest)

	return stmt
}

func (p *Parser) parseExpr(prec int) ast.Expr {
	left := p.parseUnary()

	for {
		peekPrec := p.tok.Precedence()
		if peekPrec <= prec {
			return left
		}
		left = p.parseBinary(left)
	}
}

func (p *Parser) parseUnary() ast.Expr {
	if p.tok == token.NOT || p.tok == token.MINUS {
		op := p.tok
		p.advance()
		left := p.parseExpr(token.PrecUnary)
		return &ast.UnaryExpr{Op: op, Left: left}
	}

	return p.parsePrimary()
}

func (p *Parser) parseBinary(left ast.Expr) *ast.BinaryExpr {
	op, prec := p.tok, p.tok.Precedence()
	p.advance()
	right := p.parseExpr(prec)
	return &ast.BinaryExpr{Op: op, Left: left, Right: right}
}

func (p *Parser) parsePrimary() ast.Expr {
	switch p.tok {
	case token.NUMBER:
		return p.parseNumber()
	case token.STRING:
		return p.parseString()
	case token.LPAREN:
		return p.parseGroup()
	case token.IDENTIFIER:
		return p.parseIdentifier()
	case token.TRUE:
		return p.parseTrue()
	case token.FALSE:
		return p.parseFalse()
	case token.NIL:
		return p.parseNil()
	case token.ILLEGAL:
		p.errors.Add(fmt.Sprintf("illegal token: %s", p.lit), p.file.LocationFor(p.loc))
		p.advance()
		return &ast.BadExpr{}
	default:
		p.expectError("expression")
		p.skipTo(stmtEnd)
		return &ast.BadExpr{}
	}
}

func (p *Parser) parseNumber() *ast.NumberLit {
	_, lit := p.assertToken(token.NUMBER)
	num, err := strconv.ParseFloat(lit, 64)
	assert(err == nil, fmt.Sprintf("parseNumber: %v", err))

	return &ast.NumberLit{Value: num}
}

func (p *Parser) parseString() *ast.StringLit {
	_, lit := p.assertToken(token.STRING)

	return &ast.StringLit{Value: lit}
}

func (p *Parser) parseIdentifier() *ast.IdentExpr {
	_, lit := p.assertToken(token.IDENTIFIER)

	return &ast.IdentExpr{Name: lit}
}

func (p *Parser) parseTrue() *ast.BoolExpr {
	p.assertToken(token.TRUE)

	return &ast.BoolExpr{Value: true}
}

func (p *Parser) parseFalse() *ast.BoolExpr {
	p.assertToken(token.FALSE)

	return &ast.BoolExpr{Value: false}
}

func (p *Parser) parseNil() *ast.NilExpr {
	p.assertToken(token.NIL)

	return &ast.NilExpr{}
}

func (p *Parser) parseGroup() *ast.GroupExpr {
	p.assertToken(token.LPAREN)
	expr := p.parseExpr(token.PrecLowest)
	p.expect(token.RPAREN)

	return &ast.GroupExpr{Expression: expr}
}

var stmtStart = map[token.Token]bool{
	token.NUMBER:     true,
	token.IDENTIFIER: true,
	token.STRING:     true,
	token.ILLEGAL:    true,
}

var stmtEnd = map[token.Token]bool{
	token.SEMI: true,
}

func (p *Parser) expectSemi() {
	if p.tok == token.SEMI {
		p.advance()
		return
	}
	p.expectError(token.SEMI.String())
	p.skipTo(stmtStart)
}

func (p *Parser) skipTo(to map[token.Token]bool) {
	for p.tok != token.EOF && !to[p.tok] {
		p.advance()
	}
}

func (p *Parser) expectError(expect string) {
	p.errors.Add(fmt.Sprintf("expected %q, got %q", expect, p.tok), p.file.LocationFor(p.loc))
}

func (p *Parser) expect(tok token.Token) {
	if p.tok != tok {
		p.expectError(tok.String())
	}
	p.advance()
}

func (p *Parser) assertToken(expect token.Token) (token.Token, string) {
	tok, lit := p.tok, p.lit
	assert(tok == expect, fmt.Sprintf("assertToken: expected tok %q, got %q", expect, tok))
	p.advance()

	return tok, lit
}

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func (p *Parser) advance() {
	p.tok, p.lit, p.loc = p.sc.Scan()
}
