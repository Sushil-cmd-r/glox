package lexer

import (
	"os"
	"testing"
)

func TestNext(t *testing.T) {
	tests := []string{
		"../main.glox:1:1: 12 -> Number",
		"../main.glox:1:4: + -> +",
		"../main.glox:1:6: 12 -> Number",
		"../main.glox:2:1: - -> -",
		"../main.glox:2:2: / -> /",
		"../main.glox:2:3: * -> *",
		"../main.glox:2:4: > -> >",
		"../main.glox:2:5: + -> +",
		"../main.glox:2:6: <= -> <=",
		"../main.glox:3:1: == -> ==",
		"../main.glox:3:3: 12 -> Number",
		"../main.glox:4:1: 32.53 -> Number",
		"../main.glox:4:7: >= -> >=",
		"../main.glox:4:10: 42 -> Number",
		"../main.glox:5:1: 10 -> Number",
		"../main.glox:5:4: != -> !=",
		"../main.glox:5:7: 24506343 -> Number",
		"../main.glox:7:1: ! -> !",
		"../main.glox:7:2: 2424 -> Number",
		"../main.glox:8:1: let -> Let",
		"../main.glox:8:5: five -> Identifier",
		"../main.glox:8:10: = -> =",
		"../main.glox:8:12: 5 -> Number",
		"../main.glox:8:13: ; -> ;",
	}

	content, err := os.ReadFile("../main.glox")
	if err != nil {
		panic(err)
	}

	l := New(string(content), "../main.glox")

	for i, tt := range tests {
		tok := l.Next()

		if tok.String() != tt {
			t.Fatalf("test[%d] failed: expected: %q, got: %q", i, tt, tok.String())
		}
	}
}
