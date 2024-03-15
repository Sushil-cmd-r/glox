package ast

import "fmt"

func (ie *IdentExpr) Stringify(_ int) []byte { return nil }

func (nl *NumberLiteral) Stringify(depth int) []byte {
	buf := make([]byte, 0)
	for i := 0; i < depth; i++ {
		buf = append(buf, '\t')
	}
	num := []byte(fmt.Sprintf("Number: %.2f", nl.Value))
	buf = append(buf, num...)
	buf = append(buf, '\n')

	return buf
}

func (pe *PrefixExpr) Stringify(depth int) []byte {
	buf := make([]byte, 0)
	for i := 0; i < depth; i++ {
		buf = append(buf, '\t')
	}
	op := []byte(fmt.Sprintf("Unary: %s\n", pe.Operator))
	val := pe.Right.Stringify(depth + 1)

	buf = append(buf, op...)
	buf = append(buf, val...)

	return buf
}

func (ie *InfixExpr) Stringify(depth int) []byte {
	buf := make([]byte, 0)
	for i := 0; i < depth; i++ {
		buf = append(buf, '\t')
	}
	op := []byte(fmt.Sprintf("Binary: %s\n", ie.Operator))
	buf = append(buf, op...)

	left := ie.Left.Stringify(depth + 1)
	right := ie.Right.Stringify(depth + 1)

	buf = append(buf, left...)
	buf = append(buf, right...)

	return buf
}
