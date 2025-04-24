package main

import (
	"fmt"
	"go-pascal/interpreter"
	"go-pascal/lexer"
	"go-pascal/parser"
)

func main() {
	input := "2 + 3 * 4"

	l := lexer.New(input)
	p := parser.New(l)
	expr := p.ParseExpression()
	result := interpreter.Eval(expr)

	fmt.Printf("Result: %d\n", result)
}
