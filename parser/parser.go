package parser

import (
	"fmt"

	"github.com/sushil-cmd-r/glox/ast"
	"github.com/sushil-cmd-r/glox/lexer"
	"github.com/sushil-cmd-r/glox/token"
)

type Parser struct {
	lex    *lexer.Lexer
	errors []error

	currTok token.Token
	peekTok token.Token
}

func New(lex *lexer.Lexer) *Parser {
	prs := &Parser{
		lex:    lex,
		errors: make([]error, 0),
	}

	prs.advance()
	prs.advance()

	return prs
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
		return nil, fmt.Errorf("unknown statement: %s", p.currTok.Literal)
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
