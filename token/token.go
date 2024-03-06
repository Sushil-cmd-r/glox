package token

import "github.com/sushil-cmd-r/glox/location"

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
	"fn":  Function,
	"let": Let,
}
