package token

import (
	"fmt"

	"github.com/sushil-cmd-r/glox/location"
)

type Token struct {
	Type    TokenType
	Literal string
	Loc     location.Location
}

func New(Type TokenType, Literal string, Loc location.Location) Token {
	return Token{
		Type:    Type,
		Literal: Literal,
		Loc:     Loc,
	}
}

var KeyWords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func (t Token) String() string {
	return fmt.Sprintf(
		"%s:%d:%d: %s -> %s",
		t.Loc.FilePath,
		t.Loc.Row,
		t.Loc.Col,
		t.Literal,
		t.Type,
	)
}
