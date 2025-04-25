package main

import (
	"fmt"
	"go-pascal/interpreter"
	"go-pascal/lexer"
	"go-pascal/parser"
)

func main() {
	input := `
    program Demo;

    var x: integer;

    begin
    x := 5 * 3;
    writeln(x);
    end.
    `

	// Step 1: Lexical analysis
	l := lexer.New(input)

	// Step 2: Parsing
	p := parser.New(l)
	prog := p.ParseProgram()

	// Step 3: Check for parsing errors
	if p.HasErrors() {
		fmt.Println("Parsing errors encountered:")
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		return
	}

	// Step 4: Create a new environment for interpretation
	env := interpreter.NewEnviroment()

	// Step 5: Interpret the program
	if err := interpreter.EvalProgram(prog, env); err != nil {
		fmt.Println("Runtime error encountered:")
		fmt.Println(err.Error())
		return
	}

	// Step 6: Successful execution
	fmt.Println("Program executed successfully.")
}
