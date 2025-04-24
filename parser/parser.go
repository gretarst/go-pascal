package parser

import (
	"fmt"
	"go-pascal/lexer"
	"go-pascal/token"
	"strconv"
)

type Expr interface{}

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
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseExpression() Expr {
	return p.parseAddition()
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
	fmt.Printf("Expect next token to be %s, got %s instead\n", t, p.peekToken.Type)
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
	case token.INT:
		val, _ := strconv.Atoi(p.curToken.Literal)
		lit := &IntegerLiteral{Value: val}
		p.nextToken()
		return lit
	case token.LPAREN:
		p.nextToken()
		expr := p.ParseExpression()
		p.expectPeek(token.RPAREN)
		return expr
	default:
		fmt.Printf("Unexpected token: %s\n", p.curToken.Type)
		return nil
	}

}
