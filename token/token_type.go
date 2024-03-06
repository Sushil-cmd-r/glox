package token

type TokenType string

const (
	Illegal = "Illegal"
	Eof     = "Eof"

	Identifier = "Identifier"
	Number     = "Number"

	// Operators
	Assign = "="
	Plus   = "+"

	// Delimiters
	Comma = ","
	Semi  = ";"

	LParen = "("
	RParen = ")"
	LCurly = "{"
	RCurly = "}"

	// Keywords
	Function = "fn"
	Let      = "let"
)
