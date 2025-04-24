package token

import "strings"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // main, x, y, etc.
	INT   = "INT"   // 12345

	// Operators
	ASSIGN = ":="
	PLUS   = "+"
	MINUS  = "-"
	STAR   = "*"
	SLASH  = "/"
	EQUAL  = "="
	LT     = "<"
	GT     = ">"
	LE     = "<="
	GE     = ">="
	NEQ    = "<>"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	DOT       = "."

	// Keywords
	PROGRAM = "PROGRAM"
	VAR     = "VAR"
	BEGIN   = "BEGIN"
	END     = "END"
	IF      = "IF"
	THEN    = "THEN"
	WHILE   = "WHILE"
	DO      = "DO"
	WRITELN = "WRITELN"

	// Types
	INTEGER = "INTEGER"
)

var keywords = map[string]TokenType{
	"program": PROGRAM,
	"var":     VAR,
	"begin":   BEGIN,
	"end":     END,
	"integer": INTEGER,
	"if":      IF,
	"then":    THEN,
	"while":   WHILE,
	"do":      DO,
	"writeln": WRITELN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}
