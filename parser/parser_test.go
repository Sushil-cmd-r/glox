package parser

import (
	"testing"

	"github.com/sushil-cmd-r/glox/ast"
	"github.com/sushil-cmd-r/glox/lexer"
)

func TestPrintExpr(t *testing.T) {
	input := "1 / 2 + -3"
	/*
	   Binary: +
	       Binary: /
	           Number: 1.00
	           Number: 2.00
	       Unary: -
	           Number: 3.00
	*/
	expected := "Binary: +\n\tBinary: /\n\t\tNumber: 1.00\n\t\tNumber: 2.00\n\tUnary: -\n\t\tNumber: 3.00\n"

	lex := lexer.New(input, "test.glox")
	prs := New(lex)

	program := prs.ParseProgram()
	checkParserErrors(t, prs)
	if program == nil {
		t.Fatal("ParseProgram() returned <nil>")
	}

	expr, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatal("Stmt in not an expression")
	}

	output := expr.Expression.Stringify(0)

	if string(output) != expected {
		t.Fatalf("Print Error: \nExpected: \n%s\nGot:\n%s\n", expected, string(output))
	}
}

func TestLeTStmt(t *testing.T) {
	input := `
  let x = 5;
  let y = 10;
  let foobar = 354635;
  `

	lex := lexer.New(input, "test.glox")
	prs := New(lex)

	program := prs.ParseProgram()
	checkParserErrors(t, prs)
	if program == nil {
		t.Fatal("ParseProgram() returned <nil>")
	}
	if len(program.Stmts) != 3 {
		t.Fatalf("program.Stmts does not contain 3 statements. got=%d", len(program.Stmts))
	}

	tests := []struct {
		expectedIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Stmts[i]
		if !testLetStmt(t, stmt, tt.expectedIdent) {
			return
		}
	}
}

func TestReturnStmt(t *testing.T) {
	input := `
  return 5;
  return 10;
  return 538226;
  `
	lex := lexer.New(input, "test.glox")
	prs := New(lex)

	program := prs.ParseProgram()
	checkParserErrors(t, prs)
	if program == nil {
		t.Fatal("ParseProgram() returned <nil>")
	}
	if len(program.Stmts) != 3 {
		t.Fatalf("program.Stmts does not contain 3 statements. got=%d", len(program.Stmts))
	}

	for _, stmt := range program.Stmts {
		retStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStmt. got=%T", stmt)
			continue
		}
		if retStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() not 'return', got=%q", retStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpr(t *testing.T) {
	input := `foobar;`
	lex := lexer.New(input, "test.glox")
	prs := New(lex)

	program := prs.ParseProgram()
	checkParserErrors(t, prs)
	if program == nil {
		t.Fatal("ParseProgram() returned <nil>")
	}
	if len(program.Stmts) != 1 {
		t.Fatalf("program.Stmts does not contain 1 statements. got=%d", len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt. got=%T ", program.Stmts[0])
	}

	ident, ok := stmt.Expression.(*ast.IdentExpr)
	if !ok {
		t.Fatalf("expr not *ast.IdentExpr, got=%T ", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Fatalf("ident.Value is not %s, got=%q", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral() is not %s, got=%q", "foobar", ident.TokenLiteral())
	}
}

func TestNumberLiteralExpre(t *testing.T) {
	input := `4;`
	lex := lexer.New(input, "test.glox")
	prs := New(lex)

	program := prs.ParseProgram()
	checkParserErrors(t, prs)
	if program == nil {
		t.Fatal("ParseProgram() returned <nil>")
	}
	if len(program.Stmts) != 1 {
		t.Fatalf("program.Stmts does not contain 1 statements. got=%d", len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt. got=%T ", program.Stmts[0])
	}

	literal, ok := stmt.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("expr not *ast.NumberLiteral, got=%T ", stmt.Expression)
	}
	if literal.Value != float64(4) {
		t.Fatalf("literal.Value is not %f, got=%f", float64(4), literal.Value)
	}
	if literal.TokenLiteral() != "4" {
		t.Fatalf("ident.TokenLiteral() is not %s, got=%q", "4", literal.TokenLiteral())
	}
}

func TestPrefixExpr(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		numVal   float64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		lex := lexer.New(tt.input, "test.glox")
		prs := New(lex)

		program := prs.ParseProgram()
		checkParserErrors(t, prs)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain %d statements, got=%d\n", 1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Stmt[0] is not *ast.ExpressionStmt, got %T", program.Stmts[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("stmt is not PrefixExpr, got=%T ", stmt.Expression)
		}

		if expr.Operator != tt.operator {
			t.Fatalf("expr.Operator is not `%s`, got=%s", tt.operator, expr.Operator)
		}

		if !testNumLiteral(t, expr.Right, tt.numVal) {
			return
		}
	}
}

func TestInfixExpr(t *testing.T) {
	infixTests := []struct {
		input    string
		leftVal  float64
		operator string
		rightVal float64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range infixTests {

		lex := lexer.New(tt.input, "test.glox")
		prs := New(lex)

		program := prs.ParseProgram()
		checkParserErrors(t, prs)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain %d statements, got=%d\n", 1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Stmt[0] is not *ast.ExpressionStmt, got %T", program.Stmts[0])
		}

		expr, ok := stmt.Expression.(*ast.InfixExpr)
		if !ok {
			t.Fatalf("stmt is not InfixExpr, got=%T ", stmt.Expression)
		}

		if !testNumLiteral(t, expr.Left, tt.leftVal) {
			return
		}

		if expr.Operator != tt.operator {
			t.Fatalf("expr.Operator is not `%s`, got=%s", tt.operator, expr.Operator)
		}

		if !testNumLiteral(t, expr.Right, tt.rightVal) {
			return
		}
	}
}

func testLetStmt(t *testing.T, s ast.Stmt, expected string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not `let`. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStmt)
	if !ok {
		t.Errorf("s not *ast.LetStmt. got=%T", s)
		return false
	}

	if letStmt.Name.Value != expected {
		t.Errorf("LetStmt.Name.Value not %q, got=%q", expected, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != expected {

		t.Errorf("LetStmt.Name.TokenLiteral() not %q, got=%q", expected, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testNumLiteral(t *testing.T, expr ast.Expr, value float64) bool {
	num, ok := expr.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("expr not NumberLiteral, got=%T", expr)
		return false
	}

	if num.Value != value {
		t.Errorf("num.Value not %f, got=%f ", value, num.Value)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors.", len(errs))
	for _, err := range errs {
		t.Errorf("parser error: %s", err.Error())
	}
	t.FailNow()
}
