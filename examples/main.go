package main

import (
	"fmt"
	"os"

	"github.com/sushil-cmd-r/glox/vm"
)

func main() {
	input := []byte(`
    let add = fn (a, b)  {
			let c = a + b
			print c
  	}

    let sub = fn (a, b) {
      print a - b;
    } 

    add(1 ,2 + 4)
    sub(4 ,2 + 4)
  `)
	// p := parser.New(input)
	// stmts, err := p.Parse()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// for _, stmt := range stmts {
	// 	fmt.Println(stmt)
	// }

	runCode(input)
}

func runCode(input []byte) {
	vm := vm.Init(false)

	if err := vm.Execute(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
