package parser

import (
	"fmt"

	"github.com/maffkipp/golox/lexer"
)

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{environment: NewEnvironment()}
}

func (i *Interpreter) Interpret(statements []Stmt) (hadErrors bool) {
	defer func() {
		if err := recover(); err != nil {
			if pe, ok := err.(ParseError); ok {
				pe.Report()
				hadErrors = true
			} else {
				panic(err)
			}
		}
	}()

	for _, stmt := range statements {
		i.execute(stmt)
	}
	return hadErrors
}

func (i *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) {
	i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt) {
	value := i.evaluate(stmt.Expression)
	fmt.Print(stringify(value))
}

func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) {

}

func (i *Interpreter) VisitVarStmt(stmt *VarStmt) {
	var val any
	// make sure this works!!
	if stmt.Initializer != nil {
		val = i.evaluate(stmt.Initializer)
	}
	i.environment.Define(stmt.Name.Lexeme, val)
}

func (i *Interpreter) VisitLiteralExpr(expr LiteralExpr) any {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr GroupingExpr) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr UnaryExpr) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case lexer.BANG:
		return !isTruthy(right)
	case lexer.MINUS:
		return -right.(float64)
	}

	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr BinaryExpr) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case lexer.MINUS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case lexer.SLASH:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case lexer.STAR:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	case lexer.PLUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}
		err := NewRuntimeError(expr.Operator, "operands must be numbers.")
		panic(err)
	case lexer.GREATER:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case lexer.GREATER_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case lexer.LESS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case lexer.LESS_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case lexer.BANG_EQUAL:
		return !isEqual(left, right)
	case lexer.EQUAL_EQUAL:
		return isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr VariableExpr) any {
	if val, err := i.environment.Get(expr.Name); err != nil {
		panic(err)
	} else {
		return val
	}
}

func (i *Interpreter) VisitAssignExpr(expr AssignExpr) any {
	value := i.evaluate(expr.Value)
	i.environment.Assign(expr.Name, value)
	return value
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func stringify(val any) string {
	if val == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", val)
}

func isTruthy(val any) bool {
	if val == nil {
		return false
	} else if v, ok := val.(bool); ok {
		return v
	}
	return true
}

func isEqual(left any, right any) bool {
	if left == nil && right == nil {
		return true
	} else if left == nil {
		return false
	}
	return left == right
}

func checkNumberOperands(operator lexer.Token, left any, right any) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}

	err := NewRuntimeError(operator, "operands must be numbers.")
	panic(err)
}
