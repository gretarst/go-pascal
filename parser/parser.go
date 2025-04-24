package parser

import (
	"go-pascal/lexer"
	"go-pascal/token"
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

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
