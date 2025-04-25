package parser

type Stmt any

type AssignStmt struct {
	Name  string
	Value Expr
}

type PrintStmt struct {
	Argument Expr
}

type Program struct {
	Name         string
	Declarations []Stmt
	Main         *CompoundStmt
}

type CompoundStmt struct {
	Statements []Stmt
}

type VarDecl struct {
	Name string
	Type string
}
