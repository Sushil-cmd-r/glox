package lexer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sushil-cmd-r/glox/location"
	"github.com/sushil-cmd-r/glox/token"
)

type Lexer struct {
	FileName   string
	Source     string
	Offset     int
	RdOffset   int
	LineOffset int
	Line       int
}

func New(path string) *Lexer {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	fileName := filepath.Base(path)
	fmt.Print(fileName)
	return &Lexer{
		FileName:   fileName,
		Source:     string(content),
		Offset:     0,
		RdOffset:   0,
		LineOffset: 0,
		Line:       1,
	}
}

func (l *Lexer) Next() token.Token {
	for !l.isAtEnd() {
		l.Offset = l.RdOffset
		ch := l.advance()

		loc := location.New(l.FileName, l.Line, l.LineOffset)

		switch ch {
		// Single Character Tokens
		case '+':
			return token.New(token.Plus, "+", loc)
		case '-':
			return token.New(token.Minus, "-", loc)
		case '*':
			return token.New(token.Star, "*", loc)
		case '{':
			return token.New(token.LCurly, "{", loc)
		case '}':
			return token.New(token.RCurly, "}", loc)
		case '(':
			return token.New(token.LParen, "(", loc)
		case ')':
			return token.New(token.RParen, ")", loc)
		case ',':
			return token.New(token.Comma, ",", loc)
		case ';':
			return token.New(token.Semi, ";", loc)

			// Double character Tokens
		case '=':
			if l.peek() == '=' {
				l.advance()
				return token.New(token.Equal, "==", loc)
			}
			return token.New(token.Assign, "=", loc)

		case '!':
			if l.peek() == '=' {
				l.advance()
				return token.New(token.NotEq, "!=", loc)
			}
			return token.New(token.Bang, "!", loc)

		case '>':
			if l.peek() == '=' {
				l.advance()
				return token.New(
					token.GreaterEq,
					">=",
					loc,
				)
			}
			return token.New(token.GreaterThan, ">", loc)

		case '<':
			if l.peek() == '=' {
				l.advance()
				return token.New(token.LessEq, "<=", loc)
			}
			return token.New(token.LessThan, "<", loc)

			// Ignore Comments (Only Inline Comments Supported)
		case '/':
			if l.peek() == '/' {
				for l.peek() != '\n' {
					l.advance()
				}
			} else {
				return token.New(token.Slash, "/", loc)
			}

		// whitespaces and new lines
		case ' ', '\t', '\r':
		case '\n':
			l.LineOffset = 0
			l.Line += 1

		default:
			if isNum(ch) {
				return l.number()
			} else if isChar(ch) {
				return l.Identifier()
			}
			return token.New(token.Illegal, "Unknown Character", loc)
		}
	}

	return token.New(token.Eof, "Eof", location.New(l.FileName, l.Line, l.LineOffset+1))
}

func (l *Lexer) number() token.Token {
	pos := l.LineOffset
	dotCnt := 0

	for isNum(l.peek()) || l.peek() == '.' {
		if l.peek() == '.' {
			dotCnt += 1
		}
		l.advance()
	}

	tok := l.Source[l.Offset:l.RdOffset]
	loc := location.New(l.FileName, l.Line, pos)

	if dotCnt >= 2 {
		return token.New(token.Illegal, tok, loc)
	}

	return token.New(token.Number, tok, loc)
}

func (l *Lexer) Identifier() token.Token {
	pos := location.New(l.FileName, l.Line, l.LineOffset)

	for isChar(l.peek()) {
		l.advance()
	}

	tok := l.Source[l.Offset:l.RdOffset]
	tokType, ok := token.KeyWords[tok]
	if !ok {
		return token.New(token.Identifier, tok, pos)
	}
	return token.New(tokType, tok, pos)
}

func (l *Lexer) advance() byte {
	l.RdOffset += 1
	l.LineOffset += 1
	return l.Source[l.RdOffset-1]
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return ' '
	}
	return l.Source[l.RdOffset]
}

func (l *Lexer) isAtEnd() bool {
	return l.RdOffset >= len(l.Source)
}

func isNum(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isChar(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}
