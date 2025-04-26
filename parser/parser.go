package parser

import (
	"fmt"
	"pastel/lexer"
	"pastel/token"
	"strconv"
)

type Expr any

type IntegerLiteral struct {
	Value int
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []*ParserError
}

type Identifier struct {
	Value string
}

// New creates a new Parser instance with the given lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

func (p *Parser) Errors() []*ParserError {
	return p.errors
}

// ParseExpression parses an expression in Pascal.
// Expressions include arithmetic operations like addition, subtraction, multiplication, and division.
func (p *Parser) ParseExpression() Expr {
	return p.parseAddition()
}

// ParseProgram parses a complete Pascal program.
// A Pascal program starts with the 'program' keyword, followed by declarations and a main compound statement.
func (p *Parser) ParseProgram() *Program {
	prog := &Program{}

	if p.curToken.Type == token.PROGRAM {
		// Advance to the next token after 'program' keyword
		p.nextToken()
	} else {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected 'program' keyword",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "A Pascal program must start with the 'program' keyword.",
		})
		return nil
	}

	if p.curToken.Type != token.IDENT {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected program name",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "The 'program' keyword must be followed by an identifier.",
		})
		return nil
	}

	// Advance to the next token after the program name
	prog.Name = p.curToken.Literal
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected semicolon",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Statements must end with a semicolon.",
		})
		return nil
	}

	// Advance to the next token after the semicolon
	p.nextToken()

	var decls []Stmt
	for p.curToken.Type == token.VAR {
		decl := p.parseVarDecl()
		if decl != nil {
			decls = append(decls, decl)
		}
	}
	prog.Declarations = decls

	if p.curToken.Type != token.BEGIN {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected 'begin' block",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "A Pascal program must have a 'begin' block to define its main body.",
		})
		return nil
	}

	// Parse the compound statement starting with 'begin'
	stmt := p.parseCompound()
	compound, ok := stmt.(*CompoundStmt)
	if !ok {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected compound statement",
			Detail: "The main body of the program must be a compound statement.",
			Hint:   "Ensure the program's main body starts with 'begin' and ends with 'end'.",
		})
		return nil
	}

	prog.Main = compound

	if p.curToken.Type != token.DOT {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected '.' at the end of the program",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "A Pascal program must end with a period ('.').",
		})
		return nil
	}

	return prog
}

// ParseStatement parses a single Pascal statement.
// Statements include assignments, compound statements, and print statements.
func (p *Parser) parseStatement() Stmt {
	switch p.curToken.Type {
	case token.IDENT:
		// Look ahead to see if this is an assignment (IDENT := ...)
		if p.peekToken.Type == token.ASSIGN {
			return p.parseAssignment()
		}
		p.errors = append(p.errors, &ParserError{
			Msg:    fmt.Sprintf("Unexpected identifier '%s'", p.curToken.Literal),
			Detail: "This identifier is not part of an assignment or recognized statement.",
			Hint:   "Make sure you're using ':=' for assignments or a known keyword like 'writeln'.",
		})
		p.nextToken()
		return nil

	case token.WRITELN:
		return p.parsePrint()

	case token.BEGIN:
		return p.parseCompound()

	default:
		return nil
	}
}

// ParseAssignment parses an assignment statement in Pascal.
// Assignment statements use the ':=' operator to assign values to variables.
func (p *Parser) parseAssignment() Stmt {
	name := p.curToken.Literal // We are on IDENT

	if !p.expectPeek(token.ASSIGN) {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected ':=' after identifier",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.peekToken.Literal, p.peekToken.Type),
			Hint:   "Assignments must use the ':=' operator.",
		})
		return nil
	}

	p.nextToken()
	value := p.ParseExpression()

	if !p.curTokenIs(token.SEMICOLON) {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected semicolon at the end of assignment",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Assignments must end with a semicolon.",
		})
		return nil
	}

	return &AssignStmt{Name: name, Value: value}
}

// ParseCompound parses a compound statement in Pascal.
// Compound statements start with 'begin', contain multiple statements, and end with 'end'.
func (p *Parser) parseCompound() Stmt {
	stmts := []Stmt{}

	p.nextToken()

	for p.curToken.Type != token.END && p.curToken.Type != token.EOF && p.curToken.Type != token.DOT {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}

		p.nextToken()
	}

	return &CompoundStmt{Statements: stmts}
}

// ParsePrint parses a print statement in Pascal.
// Print statements use the 'writeln' keyword to output values.
func (p *Parser) parsePrint() Stmt {
	if !p.curTokenIs(token.WRITELN) {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected 'writeln' keyword",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Use 'writeln' to print values.",
		})
		return nil
	}

	// Advance to the next token after 'writeln'
	p.nextToken()

	if !p.curTokenIs(token.LPAREN) {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected '(' after 'writeln'",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "The 'writeln' keyword must be followed by parentheses containing the argument.",
		})
		return nil
	}

	// Advance to the next token after '('
	p.nextToken()

	arg := p.ParseExpression()

	if !p.curTokenIs(token.RPAREN) {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected ')' after writeln argument",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Ensure the argument to 'writeln' is enclosed in parentheses.",
		})
		return nil
	}

	// Advance to the next token after ')'
	p.nextToken()

	if !p.curTokenIs(token.SEMICOLON) {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected ';' after writeln",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Statements must end with a semicolon.",
		})
		return nil
	}

	// Advance to the next token after the semicolon
	p.nextToken()

	return &PrintStmt{Argument: arg}
}

func PrintExpr(expr Expr, indent string) {
	switch e := expr.(type) {
	case *IntegerLiteral:
		fmt.Printf("%sInteger: %d\n", indent, e.Value)
	case *BinaryExpr:
		fmt.Printf("%sBinaryExpr: %s\n", indent, e.Operator.Literal)
		PrintExpr(e.Left, indent+"  ")
		PrintExpr(e.Right, indent+"  ")
	default:
		fmt.Printf("%sUnknown node type\n", indent)
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, &ParserError{
		Msg:    fmt.Sprintf("Expected next token to be %s", t),
		Detail: fmt.Sprintf("Got %q (%s) instead.", p.peekToken.Literal, p.peekToken.Type),
		Hint:   "Check the syntax of your program.",
	})
	return false
}

func (p *Parser) parseAddition() Expr {
	left := p.parseMultiplication()

	for p.curTokenIs(token.PLUS) || p.curTokenIs(token.MINUS) {
		op := p.curToken
		p.nextToken()
		right := p.parseMultiplication()
		left = &BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left
}

func (p *Parser) parseMultiplication() Expr {
	left := p.parsePrimary()

	for p.curTokenIs(token.STAR) || p.curTokenIs(token.SLASH) {
		op := p.curToken
		p.nextToken()
		right := p.parsePrimary()
		left = &BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left
}

func (p *Parser) parsePrimary() Expr {
	switch p.curToken.Type {
	case token.LPAREN:
		p.nextToken() // Advance from '(' to first token inside

		expr := p.ParseExpression()

		if !p.curTokenIs(token.RPAREN) {
			p.errors = append(p.errors, &ParserError{
				Msg:    "Expected closing parenthesis",
				Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
				Hint:   "Ensure all opening parentheses have matching closing parentheses.",
			})
			return nil
		}

		p.nextToken() // Consume ')'
		return expr

	case token.INT:
		val, _ := strconv.Atoi(p.curToken.Literal)
		lit := &IntegerLiteral{Value: val}
		p.nextToken()
		return lit

	case token.IDENT:
		ident := &Identifier{Value: p.curToken.Literal}
		p.nextToken()
		return ident

	default:
		p.errors = append(p.errors, &ParserError{
			Msg:    "Unexpected token in primary expression",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Check the syntax of your expression.",
		})
		return nil
	}
}

func (p *Parser) parseVarDecl() Stmt {
	// Advance to the next token after 'var'
	p.nextToken()

	if p.curToken.Type != token.IDENT {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected variable name after 'var'",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Variable declarations must start with a valid identifier.",
		})
		return nil
	}

	name := p.curToken.Literal

	// Advance to the next token after the variable name
	p.nextToken()

	if p.curToken.Type != token.COLON {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected ':' after variable name",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Variable declarations must specify a type after the colon.",
		})
		return nil
	}

	// Advance to the next token after ':'
	p.nextToken()

	if p.curToken.Type != token.INTEGER {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected 'integer' type for variable",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Currently, only 'integer' type is supported for variables.",
		})
		return nil
	}

	varType := p.curToken.Literal

	// Advance to the next token after the type
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		p.errors = append(p.errors, &ParserError{
			Msg:    "Expected ';' after variable declaration",
			Detail: fmt.Sprintf("Got %q (%s) instead.", p.curToken.Literal, p.curToken.Type),
			Hint:   "Variable declarations must end with a semicolon.",
		})
		return nil
	}

	// Advance to the next token after the semicolon
	p.nextToken()

	return &VarDecl{Name: name, Type: varType}
}
