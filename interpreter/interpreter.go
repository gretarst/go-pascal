package interpreter

import "go-pascal/parser"

func Eval(expr parser.Expr) int {
	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		return e.Value
	case *parser.BinaryExpr:
		left := Eval(e.Left)
		right := Eval(e.Right)

		switch e.Operator.Type {
		case "+":
			return left + right
		case "-":
			return left - right
		case "*":
			return left * right
		case "/":
			return left / right
		default:
			panic("Unknown operator: " + string(e.Operator.Type))
		}
	default:
		panic("Unknown expression type")
	}
}
