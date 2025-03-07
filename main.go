package main

import (
	"fmt"
	"os"

	"github.com/sushil-cmd-r/glox/vm"
)

func mxain() {
	arr := [4]int{1, 2, 3, 4}

	sarr := arr[1:]
	fmt.Println(len(arr) - len(sarr))
}

func main() {
	input := []byte(`
    function add (a, b)  {
			let c = a + b
			print c
  	}

    add(1, 2)
  `)

	runCode(input)
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
