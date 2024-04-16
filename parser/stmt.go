package parser

import "github.com/maffkipp/golox/lexer"

type Stmt interface {
	Accept(StmtVisitor)
}

type StmtVisitor interface {
	VisitBlockStmt(*BlockStmt)
	VisitExpressionStmt(*ExpressionStmt)
	VisitPrintStmt(*PrintStmt)
	VisitVarStmt(*VarStmt)
}

type BlockStmt struct {
	Statements []Stmt
}

func NewBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{Statements: statements}
}

func (b *BlockStmt) Accept(visitor StmtVisitor) {
	visitor.VisitBlockStmt(b)
}

type ExpressionStmt struct {
	Expression Expr
}

func NewExpressionStmt(expression Expr) *ExpressionStmt {
	return &ExpressionStmt{Expression: expression}
}

func (e *ExpressionStmt) Accept(visitor StmtVisitor) {
	visitor.VisitExpressionStmt(e)
}

type PrintStmt struct {
	Expression Expr
}

func NewPrintStmt(expression Expr) *PrintStmt {
	return &PrintStmt{Expression: expression}
}

func (p *PrintStmt) Accept(visitor StmtVisitor) {
	visitor.VisitPrintStmt(p)
}

type VarStmt struct {
	Name        lexer.Token
	Initializer Expr
}

func NewVarStmt(name lexer.Token, initializer Expr) *VarStmt {
	return &VarStmt{Name: name, Initializer: initializer}
}

func (v *VarStmt) Accept(visitor StmtVisitor) {
	visitor.VisitVarStmt(v)
}
