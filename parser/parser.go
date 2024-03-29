package parser

import (
	"fmt"
	"strconv"

	"github.com/sushil-cmd-r/glox/ast"
	"github.com/sushil-cmd-r/glox/lexer"
	"github.com/sushil-cmd-r/glox/token"
)

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(lhs ast.Expr) ast.Expr // argument -> left hand side of expression
)

type Precedence int

const (
	_ Precedence = iota
	LOWEST
	EQUAL       // ==
	LESSGREATER // > or <
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // -X or !X
	CALL        // fn()
)

var precedences = map[token.TokenType]Precedence{
	token.Equal:       EQUAL,
	token.NotEq:       EQUAL,
	token.LessThan:    LESSGREATER,
	token.GreaterThan: LESSGREATER,
	token.Plus:        SUM,
	token.Minus:       SUM,
	token.Slash:       PRODUCT,
	token.Star:        PRODUCT,
}

type Parser struct {
	lex    *lexer.Lexer
	errors []error

	currTok token.Token
	peekTok token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lex *lexer.Lexer) *Parser {
	prs := &Parser{
		lex:            lex,
		errors:         make([]error, 0),
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	prs.registerPrefixFns(token.Identifier, prs.parseIdentifier)
	prs.registerPrefixFns(token.Number, prs.parserNumberLiteral)
	prs.registerPrefixFns(token.Bang, prs.parsePrefixExpr)
	prs.registerPrefixFns(token.Minus, prs.parsePrefixExpr)

	prs.registerInfixFns(token.Plus, prs.parseInfixExpr)
	prs.registerInfixFns(token.Minus, prs.parseInfixExpr)
	prs.registerInfixFns(token.Slash, prs.parseInfixExpr)
	prs.registerInfixFns(token.Star, prs.parseInfixExpr)
	prs.registerInfixFns(token.LessThan, prs.parseInfixExpr)
	prs.registerInfixFns(token.GreaterThan, prs.parseInfixExpr)
	prs.registerInfixFns(token.Equal, prs.parseInfixExpr)
	prs.registerInfixFns(token.NotEq, prs.parseInfixExpr)

	prs.advance()
	prs.advance()

	return prs
}

func (p *Parser) registerPrefixFns(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfixFns(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) peekPrecedence() Precedence {
	if p, ok := precedences[p.peekTok.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecedence() Precedence {
	if p, ok := precedences[p.currTok.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Stmts = make([]ast.Stmt, 0)

	for p.currTok.Type != token.Eof {
		stmt, err := p.parseStmt()
		if err != nil {
			p.errors = append(p.errors, err)
			// Synchronize until next statement start
			p.synchronize()
		} else {
			program.Stmts = append(program.Stmts, stmt)
		}

		p.advance()
	}

	return program
}

func (p *Parser) parseStmt() (ast.Stmt, error) {
	switch p.currTok.Type {
	case token.Let:
		return p.parseLetStmt()
	case token.Return:
		return p.parseReturnStmt()
	default:
		return p.parseExpressionStmt()
	}
}

func (p *Parser) parseLetStmt() (*ast.LetStmt, error) {
	stmt := &ast.LetStmt{
		Token: p.currTok,
	}

	if err := p.expect(token.Identifier); err != nil {
		return nil, err
	}

	stmt.Name = &ast.IdentExpr{
		Token: p.currTok,
		Value: p.currTok.Literal,
	}

	if err := p.expect(token.Assign); err != nil {
		return nil, err
	}

	// TODO: Skip Value for Now
	for !p.currTokIs(token.Semi) {
		p.advance()
	}

	return stmt, nil
}

func (p *Parser) parseReturnStmt() (*ast.ReturnStmt, error) {
	stmt := &ast.ReturnStmt{
		Token: p.currTok,
	}
	p.advance()

	// TODO: Skip return value for now
	for !p.currTokIs(token.Semi) {
		p.advance()
	}

	return stmt, nil
}

func (p *Parser) parseExpressionStmt() (*ast.ExpressionStmt, error) {
	stmt := &ast.ExpressionStmt{
		Token: p.currTok,
	}

	stmt.Expression = p.parseExpression(LOWEST)
	if stmt.Expression == nil {
		return nil, fmt.Errorf(`unknown expression %s `, p.currTok.Literal)
	}

	if p.peekTokIs(token.Semi) {
		p.advance()
	}

	return stmt, nil
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expr {
	prefix := p.prefixParseFns[p.currTok.Type]
	if prefix == nil {
		return nil
	}
	leftExpr := prefix()

	for !p.peekTokIs(token.Semi) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekTok.Type]
		if infix == nil {
			return leftExpr
		}

		p.advance()
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.IdentExpr{
		Token: p.currTok,
		Value: p.currTok.Literal,
	}
}

func (p *Parser) parserNumberLiteral() ast.Expr {
	lit := &ast.NumberLiteral{
		Token: p.currTok,
	}

	value, err := strconv.ParseFloat(p.currTok.Literal, 64)
	if err != nil {
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpr() ast.Expr {
	expr := &ast.PrefixExpr{
		Token:    p.currTok,
		Operator: p.currTok.Literal,
	}
	p.advance()

	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{
		Token:    p.currTok,
		Operator: p.currTok.Literal,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.advance()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) synchronize() {
	for p.peekTok.Type != token.Eof {
		switch p.peekTok.Type {
		case token.Let, token.Return:
			return
		}

		p.advance()
	}
}

// Errors
func (p *Parser) unexpectedTokError(expected token.TokenType) error {
	return fmt.Errorf("expected next token to be %s, got %s instead.", expected, p.peekTok.Type)
}

func (p *Parser) advance() {
	p.currTok = p.peekTok
	p.peekTok = p.lex.Next()
}

func (p *Parser) currTokIs(tt token.TokenType) bool {
	return p.currTok.Type == tt
}

func (p *Parser) peekTokIs(tt token.TokenType) bool {
	return p.peekTok.Type == tt
}

func (p *Parser) expect(tt token.TokenType) error {
	if p.peekTokIs(tt) {
		p.advance()
		return nil
	}

	return p.unexpectedTokError(tt)
}
