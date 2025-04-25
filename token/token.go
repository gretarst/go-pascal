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
	IDENT = "IDENT" // e.g., variable names
	INT   = "INT"   // e.g., 123

	// Operators
	ASSIGN = "ASSIGN" // :=
	PLUS   = "PLUS"   // +
	MINUS  = "MINUS"  // -
	STAR   = "STAR"   // *
	SLASH  = "SLASH"  // /

	EQUAL = "EQUAL" // =
	LT    = "LT"    // <
	GT    = "GT"    // >
	LE    = "LE"    // <=
	GE    = "GE"    // >=
	NEQ   = "NEQ"   // <>

	// Delimiters
	COMMA     = "COMMA"     // ,
	SEMICOLON = "SEMICOLON" // ;
	COLON     = "COLON"     // :
	LPAREN    = "LPAREN"    // (
	RPAREN    = "RPAREN"    // )
	DOT       = "DOT"       // .

	// Keywords
	AND       = "AND"
	ARRAY     = "ARRAY"
	BEGIN     = "BEGIN"
	CASE      = "CASE"
	CONST     = "CONST"
	DIV       = "DIV"
	DO        = "DO"
	DOWNTO    = "DOWNTO"
	ELSE      = "ELSE"
	END       = "END"
	FILE      = "FILE"
	FOR       = "FOR"
	FORWARD   = "FORWARD"
	FUNCTION  = "FUNCTION"
	GOTO      = "GOTO"
	IF        = "IF"
	IN        = "IN"
	LABEL     = "LABEL"
	MOD       = "MOD"
	NIL       = "NIL"
	NOT       = "NOT"
	OF        = "OF"
	OR        = "OR"
	PACKED    = "PACKED"
	PROCEDURE = "PROCEDURE"
	PROGRAM   = "PROGRAM"
	RECORD    = "RECORD"
	REPEAT    = "REPEAT"
	SET       = "SET"
	THEN      = "THEN"
	TO        = "TO"
	TYPE      = "TYPE"
	UNTIL     = "UNTIL"
	VAR       = "VAR"
	WHILE     = "WHILE"
	WITH      = "WITH"
	WRITELN   = "WRITELN"

	// Types
	INTEGER = "INTEGER"
)

var keywords = map[string]TokenType{
	"and":       AND,
	"array":     ARRAY,
	"begin":     BEGIN,
	"case":      CASE,
	"const":     CONST,
	"div":       DIV,
	"do":        DO,
	"downto":    DOWNTO,
	"else":      ELSE,
	"end":       END,
	"file":      FILE,
	"for":       FOR,
	"forward":   FORWARD,
	"function":  FUNCTION,
	"goto":      GOTO,
	"if":        IF,
	"in":        IN,
	"label":     LABEL,
	"mod":       MOD,
	"nil":       NIL,
	"not":       NOT,
	"of":        OF,
	"or":        OR,
	"packed":    PACKED,
	"procedure": PROCEDURE,
	"program":   PROGRAM,
	"record":    RECORD,
	"repeat":    REPEAT,
	"set":       SET,
	"then":      THEN,
	"to":        TO,
	"type":      TYPE,
	"until":     UNTIL,
	"var":       VAR,
	"while":     WHILE,
	"with":      WITH,
	"writeln":   WRITELN,
	"integer":   INTEGER,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}
