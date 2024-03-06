package lexer

import (
	"fmt"
	"testing"

	"github.com/sushil-cmd-r/glox/token"
)

func TestNext(t *testing.T) {
	input := "  =+{}(),;\n1.23 "

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedPos     string
	}{
		{expectedType: token.Assign, expectedLiteral: "=", expectedPos: "1:3"},
		{expectedType: token.Plus, expectedLiteral: "+", expectedPos: "1:4"},
		{expectedType: token.LCurly, expectedLiteral: "{", expectedPos: "1:5"},
		{expectedType: token.RCurly, expectedLiteral: "}", expectedPos: "1:6"},
		{expectedType: token.LParen, expectedLiteral: "(", expectedPos: "1:7"},
		{expectedType: token.RParen, expectedLiteral: ")", expectedPos: "1:8"},
		{expectedType: token.Comma, expectedLiteral: ",", expectedPos: "1:9"},
		{expectedType: token.Semi, expectedLiteral: ";", expectedPos: "1:10"},
		{expectedType: token.Number, expectedLiteral: "1.23", expectedPos: "2:1"},
		{expectedType: token.Eof, expectedLiteral: "Eof", expectedPos: "2:6"},
	}

	lex := New(input, "test")

	for i, tt := range tests {
		tok := lex.Next()

		if tok.Type != tt.expectedType {
			t.Fatalf(
				"tests[%d] - tokenType wrong, expected: %q, got: %q ",
				i,
				tt.expectedType,
				tok.Type,
			)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] - tokenLiteral wrong, expected: %q, got: %q ",
				i,
				tt.expectedLiteral,
				tok.Literal,
			)
		}

		pos := fmt.Sprintf("%d:%d", tok.Loc.Row, tok.Loc.Col)
		if pos != tt.expectedPos {
			t.Fatalf(
				"tests[%d] - tokenPos wrong, expected: %q, got: %q ",
				i,
				tt.expectedPos,
				pos,
			)
		}
	}
}
