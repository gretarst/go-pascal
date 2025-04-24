package main

import (
	"go-pascal/lexer"
	"go-pascal/parser"
)

func main() {
	input := "2 + 3 * 4"

	l := lexer.New(input)
	p := parser.New(l)
	expr := p.ParseExpression()
	parser.PrintExpr(expr, "")
}
