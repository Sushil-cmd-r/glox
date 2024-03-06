package lexer

import (
	"github.com/sushil-cmd-r/glox/location"
	"github.com/sushil-cmd-r/glox/token"
)

type Lexer struct {
	FilePath string
	Source   string
	Offset   int
	Loffset  int
	TStart   int
	Line     int
}

func New(source, filepath string) *Lexer {
	return &Lexer{
		FilePath: filepath,
		Source:   source,
		Offset:   0,
		Loffset:  0,
		TStart:   0,
		Line:     1,
	}
}

func (l *Lexer) Next() token.Token {
	for !l.isAtEnd() {
		l.TStart = l.Offset
		ch := l.advance()

		switch ch {
		case '+':
			return token.New(token.Plus, "+", location.New(l.FilePath, l.Line, l.Loffset))
		case '=':
			return token.New(token.Assign, "=", location.New(l.FilePath, l.Line, l.Loffset))
		case '{':
			return token.New(token.LCurly, "{", location.New(l.FilePath, l.Line, l.Loffset))
		case '}':
			return token.New(token.RCurly, "}", location.New(l.FilePath, l.Line, l.Loffset))
		case '(':
			return token.New(token.LParen, "(", location.New(l.FilePath, l.Line, l.Loffset))
		case ')':
			return token.New(token.RParen, ")", location.New(l.FilePath, l.Line, l.Loffset))
		case ',':
			return token.New(token.Comma, ",", location.New(l.FilePath, l.Line, l.Loffset))
		case ';':
			return token.New(token.Semi, ";", location.New(l.FilePath, l.Line, l.Loffset))
		case ' ', '\t', '\r':
			// ignore whitespaces
		case '\n':
			l.Loffset = 0
			l.Line += 1
		default:
			if isNum(ch) {
				return l.number()
			}
			return token.New(token.Illegal, "Illegal", location.New(l.FilePath, l.Line, l.Loffset))
		}
	}

	return token.New(token.Eof, "Eof", location.New(l.FilePath, l.Line, l.Loffset+1))
}

func (l *Lexer) number() token.Token {
	pos := l.Loffset
	dotCnt := 0

	for isNum(l.peek()) || l.peek() == '.' {
		if l.peek() == '.' {
			dotCnt += 1
		}
		l.advance()
	}

	tok := l.Source[l.TStart:l.Offset]
	loc := location.New(l.FilePath, l.Line, pos)

	if dotCnt >= 2 {
		return token.New(token.Illegal, "Illegal", loc)
	}

	return token.New(token.Number, tok, loc)
}

func (l *Lexer) advance() byte {
	l.Offset += 1
	l.Loffset += 1
	return l.Source[l.Offset-1]
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return ' '
	}
	return l.Source[l.Offset]
}

func (l *Lexer) isAtEnd() bool {
	return l.Offset >= len(l.Source)
}

func isNum(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
