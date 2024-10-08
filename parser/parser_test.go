package parser

import (
	"fmt"
	"testing"

	"github.com/sushil-cmd-r/glox/token"
)

func makeParser(input string) *Parser {
	f := token.NewFile("test.glox")
	return New(f, input)
}

func TestParse(t *testing.T) {
	input := `1 + 2; x
            true; nil
           "hello";`
	expect := []string{
		"(+ 1 2);\n",
		"x;\n",
		"true;\n",
		"<nil>;\n",
		"hello;\n",
	}

	p := makeParser(input)
	stmts := p.Parse()

	if p.errors.Len() != 0 {
		t.Fatal(stmts, p.errors.Error())
	}

	for i, tc := range expect {
		stmt := stmts[i]
		stmtStr := fmt.Sprintf("%s", stmt)

		if stmtStr != tc {
			t.Fatalf("TestParse [%d]: expected %s, got %s", i+1, tc, stmtStr)
		}
	}
}

func TestParseExpr(t *testing.T) {
	input := `( -x + 3.14) / "hello"`

	expect := "(/ (+ (- x) 3.14) hello)"

	p := makeParser(input)
	expr := p.parseExpr(token.PrecLowest)
	if expr == nil {
		t.Fatal("TestParseExpr: expected expression")
	}

	if p.errors.Len() > 0 {
		t.Fatal(p.errors)
	}

	exprStr := fmt.Sprintf("%s", expr)

	if expect != exprStr {
		t.Fatalf("TestParseExpr: expected expr %s, got %s", expect, exprStr)
	}
}
