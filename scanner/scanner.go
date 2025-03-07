package scanner

import "github.com/sushil-cmd-r/glox/token"

var eof byte = 0

type Scanner struct {
	source   []byte
	rdOffset int

	ch     byte
	offset int

	insertSemi bool
}

func Init(source []byte) *Scanner {
	s := &Scanner{
		source:   source,
		rdOffset: 0,

		ch:     ' ',
		offset: 0,

		insertSemi: false,
	}
	s.advance()

	return s
}

func scanToken(tok token.Token) (token.Token, string) {
	return tok, tok.String()
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || (!s.insertSemi && s.ch == '\n') || (!s.insertSemi && s.ch == '\r') {
		s.advance()
	}
}

func (s *Scanner) Scan() (tok token.Token, lit string) {
	s.skipWhitespace()
	ch := s.ch
	s.advance()

	insertSemi := false
	switch ch {
	case eof:
		if s.insertSemi {
			tok, lit = scanToken(token.SEMI)
		} else {
			tok, lit = scanToken(token.EOF)
		}
	case '+':
		tok, lit = scanToken(token.PLUS)
	case '-':
		tok, lit = scanToken(token.MINUS)
	case '*':
		tok, lit = scanToken(token.STAR)
	case '/':
		tok, lit = scanToken(token.SLASH)

	case '(':
		tok, lit = scanToken(token.LPAREN)
	case ')':
		insertSemi = true
		tok, lit = scanToken(token.RPAREN)
	case '{':
		tok, lit = scanToken(token.LCURLY)
	case '}':
		insertSemi = true
		tok, lit = scanToken(token.RCURLY)
	case ',':
		tok, lit = scanToken(token.COMMA)
	case ';':
		tok, lit = scanToken(token.SEMI)

	case '=':
		tok, lit = s.switch0(token.ASSIGN, token.EQL)
	case '!':
		tok, lit = s.switch0(token.NOT, token.NEQ)
	case '>':
		tok, lit = s.switch0(token.GTR, token.GEQ)
	case '<':
		tok, lit = s.switch0(token.LSS, token.LEQ)

	case '"':
		insertSemi = true
		tok, lit = s.scanString()

	case '\n', '\r':
		tok, lit = scanToken(token.SEMI)

	default:
		if isNum(ch) {
			tok, lit = s.scanNumber()
			insertSemi = true

		} else if isChar(ch) {
			tok, lit = s.scanIdentifier()
			if tok == token.IDENTIFIER {
				insertSemi = true
			}

		} else {
			tok, lit = token.ILLEGAL, string(ch)
			insertSemi = true
		}
	}

	s.insertSemi = insertSemi
	return
}

func (s *Scanner) scanString() (token.Token, string) {
	st := s.offset

	for s.ch != '"' {
		s.advance()
		if s.ch == eof {
			break
		}
	}

	if s.ch == eof {
		return token.ILLEGAL, "unterminated string"
	}
	lit := string(s.source[st:s.offset])
	s.advance()

	return token.STRING, lit
}

func (s *Scanner) scanNumber() (token.Token, string) {
	st := s.offset - 1
	valid := true
	dotCnt := 0

	for isNum(s.ch) || isChar(s.ch) || s.ch == '.' {
		if valid && isChar(s.ch) {
			valid = false
		}

		if s.ch == '.' {
			dotCnt += 1
		}
		s.advance()
	}

	lit := string(s.source[st:s.offset])

	if !valid || dotCnt >= 2 {
		return token.ILLEGAL, lit
	}
	return token.NUMBER, lit
}

func (s *Scanner) scanIdentifier() (token.Token, string) {
	st := s.offset - 1
	for isChar(s.ch) || isNum(s.ch) {
		s.advance()
	}

	lit := string(s.source[st:s.offset])
	return token.LookUp(lit)
}

func (s *Scanner) switch0(t1, t2 token.Token) (token.Token, string) {
	if s.ch == '=' {
		s.advance()
		return scanToken(t2)
	}
	return scanToken(t1)
}

func (s *Scanner) advance() {
	if s.atEnd() {
		s.ch = eof
		s.offset = len(s.source)
		return
	}
	s.offset = s.rdOffset
	s.ch = s.source[s.offset]
	s.rdOffset += 1
}

func (s *Scanner) atEnd() bool {
	return s.rdOffset >= len(s.source)
}

func isNum(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}
