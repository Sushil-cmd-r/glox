package main

import (
	"fmt"
	"os"

	"github.com/sushil-cmd-r/glox/parser"
	"github.com/sushil-cmd-r/glox/vm"
)

func main() {
	input := []byte(`
    function add (a, b, c)  {
			let c = a + b
			print c
  	}

    add(1 ,2 + 4)
  `)

	p := parser.New(input)
	stmts, err := p.Parse()
	if err != nil {
		fmt.Println(err)
	}

	for _, stmt := range stmts {
		fmt.Println(stmt)
	}

	// runCode(input)
}

func runCode(input []byte) {
	vm := vm.Init()

	if err := vm.Execute(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// vm.PrintCode()

	// vm.Interpret()
}
