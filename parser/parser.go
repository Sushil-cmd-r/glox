package parser

import (
	"fmt"
	"strconv"

	"github.com/sushil-cmd-r/glox/ast"
	"github.com/sushil-cmd-r/glox/scanner"
	"github.com/sushil-cmd-r/glox/token"
)

type Parser struct {
	sc *scanner.Scanner

	err *ErrorList

	tok token.Token
	lit string
}

func New(source []byte) *Parser {
	sc := scanner.Init(source)
	p := &Parser{sc: sc, err: &ErrorList{}}

	p.advance()
	return p
}

func (p *Parser) Parse() (prog []ast.Stmt, err error) {
	for p.tok != token.EOF {
		stmt := p.parseStmt()
		p.expectSemi()
		prog = append(prog, stmt)
	}

	err = p.err.Err()
	return
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.tok {
	case token.LET:
		return p.parseLetStmt()
	case token.LCURLY:
		return p.parseBlockStmt()
	case token.PRINT:
		return p.parsePrintStmt()
	case token.FUNCTION:
		return p.parseFuncStmt()
	default:
		return p.parsePriamryStmt()
	}
}

func (p *Parser) parseFuncStmt() *ast.FuncStmt {
	p.advance()
	name := p.parseIdentifier()

	var params []*ast.IdentExpr

	p.expect(token.LPAREN)
	for p.tok != token.RPAREN && p.tok != token.EOF {
		expr := p.parseIdentifier()
		params = append(params, expr)
		if p.tok == token.RPAREN {
			break
		}
		p.expect(token.COMMA)
	}

	p.expect(token.RPAREN)

	body := p.parseBlockStmt()
	return &ast.FuncStmt{Name: name, Params: params, Body: body}
}

func (p *Parser) parsePrintStmt() *ast.PrintStmt {
	p.advance()
	expr := p.parseExpr(token.PrecLowest)
	return &ast.PrintStmt{Expr: expr}
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	p.advance()
	var stmts []ast.Stmt
	for p.tok != token.EOF && p.tok != token.RCURLY {
		stmt := p.parseStmt()
		p.expectSemi()
		stmts = append(stmts, stmt)
	}

	p.expect(token.RCURLY)
	return &ast.BlockStmt{Stmts: stmts}
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	p.advance()
	name := p.parseIdentifier()

	if p.tok != token.ASSIGN {
		return &ast.LetStmt{Name: name, Value: &ast.NilExpr{}}
	}
	p.advance()
	expr := p.parseExpr(token.PrecLowest)

	return &ast.LetStmt{Name: name, Value: expr}
}

func (p *Parser) parsePriamryStmt() ast.Stmt {
	expression := p.parseExpr(token.PrecLowest)

	if p.tok == token.ASSIGN {
		p.advance()
		value := p.parseExpr(token.PrecLowest)
		return &ast.AssignStmt{Name: expression, Value: value}
	}

	if p.tok == token.LPAREN {
		var args []ast.Expr
		p.advance()
		if p.tok != token.RPAREN {
			expr := p.parseExpr(token.PrecLowest)
			args = append(args, expr)

			for p.tok == token.COMMA && p.tok != token.EOF {
				p.advance()
				expr := p.parseExpr(token.PrecLowest)
				args = append(args, expr)
			}
		}

		p.expect(token.RPAREN)
		expression = &ast.CallExpr{Callee: expression, Args: args}
	}

	return &ast.ExprStmt{Expression: expression}
}

func (p *Parser) parseExpr(prec int) ast.Expr {
	left := p.parseUnary()

	for {
		peekPrec := p.tok.Precedence()
		if peekPrec <= prec {
			break
		}
		left = p.parseBinary(left)
	}

	return left
}

func (p *Parser) parseUnary() ast.Expr {
	if p.tok == token.MINUS || p.tok == token.NOT {
		op := p.tok
		p.advance()
		left := p.parseExpr(token.PrecUnary)
		return &ast.UnaryExpr{Op: op, Left: left}
	}

	return p.parsePrimary()
}

func (p *Parser) parseBinary(left ast.Expr) ast.Expr {
	op := p.tok
	p.advance()

	right := p.parseExpr(token.PrecLowest)
	return &ast.BinaryExpr{Op: op, Left: left, Right: right}
}

func (p *Parser) parsePrimary() ast.Expr {
	switch p.tok {
	case token.NUMBER:
		return p.parseNumber()
	case token.STRING:
		return p.parseString()
	case token.IDENTIFIER:
		return p.parseIdentifier()
	case token.LPAREN:
		return p.parseGroup()

	default:
		p.expectError("expression")
		return nil
	}
}

func (p *Parser) parseNumber() *ast.NumberLit {
	_, lit := p.assertTok(token.NUMBER)

	num, err := strconv.ParseFloat(lit, 64)
	assert(err == nil, "invalid number")

	return &ast.NumberLit{Value: num}
}

func (p *Parser) parseString() *ast.StringLit {
	_, lit := p.assertTok(token.STRING)
	return &ast.StringLit{Value: lit}
}

func (p *Parser) parseIdentifier() *ast.IdentExpr {
	name := "_"
	if p.tok != token.IDENTIFIER {
		p.expectError("identifier")
	} else {
		name = p.lit
		p.advance()
	}

	return &ast.IdentExpr{Name: name}
}

func (p *Parser) parseGroup() *ast.GroupExpr {
	p.assertTok(token.LPAREN)

	expr := p.parseExpr(token.PrecLowest)
	p.expect(token.RPAREN)

	return &ast.GroupExpr{Expression: expr}
}

func (p *Parser) expectSemi() {
	p.expect(token.SEMI)
}

func (p *Parser) expect(tok token.Token) {
	if p.tok != tok {
		p.expectError("'" + tok.String() + "'")
	}
	p.advance()
	return
}

func (p *Parser) expectError(msg string) {
	msg = "expected " + msg

	switch p.tok {
	case token.NUMBER, token.IDENTIFIER:
		msg += "got, " + p.lit

	default:
		msg += ", got " + p.tok.String()
	}

	p.errors(msg)
}

func (p *Parser) errors(msg string) {
	n := p.err.Len()
	assert(n < 10, "too many errors")

	p.err.Add(msg)
}

func (p *Parser) assertTok(expect token.Token) (token.Token, string) {
	tok, lit := p.tok, p.lit

	assert(tok == expect, fmt.Sprintf("expected tok %s, got %s", expect, tok))
	p.advance()
	return tok, lit
}

func assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}

func (p *Parser) advance() {
	p.tok, p.lit = p.sc.Scan()
}
