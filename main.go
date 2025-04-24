package main

import (
	"fmt"
	"go-pascal/lexer"
	"go-pascal/token"
)

func main() {
	input := `
program Demo;
var x: integer;
var counter1: integer;
BeGiN
  counter1 := 99;
  x := 42 + 3;
  writeln(x);
end.
`

	l := lexer.New(input)

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
}
