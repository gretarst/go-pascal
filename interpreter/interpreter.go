package interpreter

import (
	"fmt"
	"go-pascal/parser"
	"go-pascal/token"
)

// EvalProgram evaluates the entire Pascal program.
func EvalProgram(prog *parser.Program, env *Environment) error {
	for _, decl := range prog.Declarations {
		if v, ok := decl.(*parser.VarDecl); ok {
			env.Set(v.Name, 0) // default value is 0
		}
	}

	if err := EvalStmt(prog.Main, env); err != nil {
		return err
	}

	return nil
}

// EvalStmt evaluates a single statement.
func EvalStmt(stmt parser.Stmt, env *Environment) error {
	switch s := stmt.(type) {
	case *parser.AssignStmt:
		if !env.Exists(s.Name) {
			return &PascalError{
				Msg:    fmt.Sprintf("Undeclared variable '%s'", s.Name),
				Detail: "This variable is being used but was never declared with a type.",
				Hint:   fmt.Sprintf("Try adding `var %s: integer;` at the top of your program.", s.Name),
			}
		}

		val, err := EvalExpr(s.Value, env)
		if err != nil {
			return err
		}
		env.Set(s.Name, val)

	case *parser.CompoundStmt:
		for _, stmt := range s.Statements {
			if err := EvalStmt(stmt, env); err != nil {
				return err
			}
		}

	case *parser.PrintStmt:
		val, err := EvalExpr(s.Argument, env)
		if err != nil {
			return err
		}
		fmt.Println(val)

	default:
		return &PascalError{
			Msg:    "Unknown statement type",
			Detail: fmt.Sprintf("Encountered an unsupported statement: %T", stmt),
			Hint:   "Ensure all statements are valid Pascal constructs.",
		}
	}

	return nil
}

// EvalExpr evaluates an expression and returns its integer value.
func EvalExpr(expr parser.Expr, env *Environment) (int, error) {
	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		return e.Value, nil

	case *parser.BinaryExpr:
		left, err := EvalExpr(e.Left, env)
		if err != nil {
			return 0, err
		}

		right, err := EvalExpr(e.Right, env)
		if err != nil {
			return 0, err
		}

		switch e.Operator.Type {
		case token.PLUS:
			return left + right, nil
		case token.MINUS:
			return left - right, nil
		case token.STAR:
			return left * right, nil
		case token.SLASH:
			if right == 0 {
				return 0, &PascalError{
					Msg:    "Division by zero",
					Detail: "An attempt was made to divide by zero.",
					Hint:   "Ensure the divisor is not zero before performing division.",
				}
			}
			return left / right, nil
		default:
			return 0, &PascalError{
				Msg:    "Unknown operator",
				Detail: fmt.Sprintf("Operator '%s' is not supported.", e.Operator.Literal),
				Hint:   "Use valid operators such as +, -, *, or /.",
			}
		}

	case *parser.Identifier:
		val, ok := env.Get(e.Value)
		if !ok {
			return 0, &PascalError{
				Msg:    fmt.Sprintf("Undefined variable '%s'", e.Value),
				Detail: "This variable is being used but was never declared or assigned a value.",
				Hint:   fmt.Sprintf("Declare the variable using `var %s: integer;` and assign it a value before use.", e.Value),
			}
		}
		return val, nil

	default:
		return 0, &PascalError{
			Msg:    "Unknown expression type",
			Detail: fmt.Sprintf("Encountered an unsupported expression: %T", expr),
			Hint:   "Ensure all expressions are valid Pascal constructs.",
		}
	}
}
