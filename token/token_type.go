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
	Minus  = "-"
	Bang   = "!"
	Star   = "*"
	Slash  = "/"

	LessThan    = "<"
	LessEq      = "<="
	GreaterThan = ">"
	GreaterEq   = ">="
	Equal       = "=="
	NotEq       = "!="

	// Delimiters
	Comma = ","
	Semi  = ";"

	LParen = "("
	RParen = ")"
	LCurly = "{"
	RCurly = "}"

	// Keywords
	Function = "Function"
	Let      = "Let"
	True     = "True"
	False    = "False"
	If       = "If"
	Else     = "Else"
	Return   = "Return"
)
