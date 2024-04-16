package parser

import "github.com/maffkipp/golox/lexer"

type Expr interface {
	Accept(ExprVisitor) any
}

type ExprVisitor interface {
	VisitUnaryExpr(UnaryExpr) any
	VisitBinaryExpr(BinaryExpr) any
	VisitGroupingExpr(GroupingExpr) any
	VisitLiteralExpr(LiteralExpr) any
	VisitVariableExpr(VariableExpr) any
	VisitAssignExpr(AssignExpr) any
}

type UnaryExpr struct {
	Operator lexer.Token
	Right    Expr
}

func NewUnaryExpr(operator lexer.Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator, right}
}

func (u UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(u)
}

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func NewBinaryExpr(left Expr, operator lexer.Token, right Expr) *BinaryExpr {
	return &BinaryExpr{left, operator, right}
}

func (b BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expression Expr
}

func NewGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{expression}
}

func (g GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func NewLiteralExpr(value any) *LiteralExpr {
	return &LiteralExpr{value}
}

func (l LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(l)
}

type VariableExpr struct {
	Name lexer.Token
}

func NewVariableExpr(name lexer.Token) *VariableExpr {
	return &VariableExpr{Name: name}
}

func (v VariableExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpr(v)
}

type AssignExpr struct {
	Name  lexer.Token
	Value Expr
}

func NewAssignExpr(name lexer.Token, value Expr) *AssignExpr {
	return &AssignExpr{Name: name, Value: value}
}

func (a AssignExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignExpr(a)
}
