package token

type Token int

const (
	ILLEGAL Token = iota // illegal
	EOF                  // eof

	NUMBER     // number
	STRING     // string
	IDENTIFIER // identifier

	ASSIGN // =
	LPAREN // (
	RPAREN // )
	LCURLY // {
	RCURLY // }
	COMMA  // ,
	SEMI   // ;
	NOT    // !

	PLUS  // +
	MINUS // -
	STAR  // *
	SLASH // /
	EQL   // ==
	NEQ   // !=
	GTR   // >
	GEQ   // >=
	LSS   // <
	LEQ   // <=

	keywordStart
	LET      // let
	NIL      // nil
	PRINT    // print
	FUNCTION // function
	FN       // fn
	keywordEnd
)

var tokens = [...]string{
	ILLEGAL: "illegal",
	EOF:     "eof",

	NUMBER:     "number",
	STRING:     "string",
	IDENTIFIER: "identifier",

	ASSIGN: "=",
	LPAREN: "(",
	RPAREN: ")",
	LCURLY: "{",
	RCURLY: "}",
	COMMA:  ",",
	SEMI:   ";",
	NOT:    "!",

	PLUS:  "+",
	MINUS: "-",
	STAR:  "*",
	SLASH: "/",
	EQL:   "==",
	NEQ:   "!=",
	GTR:   ">",
	GEQ:   ">=",
	LSS:   "<",
	LEQ:   "<=",

	LET:      "let",
	NIL:      "nil",
	PRINT:    "print",
	FUNCTION: "function",
	FN:       "fn",
}

func (tok Token) String() string {
	return tokens[tok]
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keywordStart + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = i
	}
}

func LookUp(lit string) (Token, string) {
	if keyword, ok := keywords[lit]; ok {
		return keyword, lit
	}

	return IDENTIFIER, lit
}

const (
	PrecLowest = iota
	PrecEquality
	PrecComparision
	PrecTerm
	PrecFactor
	PrecUnary
)

func (tok Token) Precedence() int {
	switch tok {
	case EQL, NEQ:
		return PrecEquality
	case GTR, LSS, GEQ, LEQ:
		return PrecComparision
	case PLUS, MINUS:
		return PrecTerm
	case STAR, SLASH:
		return PrecFactor
	default:
		return PrecLowest
	}
}
