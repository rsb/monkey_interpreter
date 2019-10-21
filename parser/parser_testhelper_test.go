package parser_test

import (
	"fmt"
	"testing"

	"github.com/rsb/monkey_interpreter/ast"
	"github.com/rsb/monkey_interpreter/parser"
	"github.com/stretchr/testify/assert"
)

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, op string, right interface{}) {
	opExpr, ok := expr.(*ast.InfixExpression)
	assert.True(t, ok, "expr is not *ast.InfixExpression, got=%T", expr)
	testLiteralExpression(t, opExpr.Left, left)

	assert.Equal(t, op, opExpr.Operator, "expr.Operator is not '%s', got=%q", op, opExpr.Operator)
	testLiteralExpression(t, opExpr.Right, right)
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, expr, int64(v))
	case int64:
		testIntegerLiteral(t, expr, v)
	case string:
		testIdentifier(t, expr, v)
	default:
		t.Errorf("type of expr not handled. got=%T", expr)
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	integ, ok := il.(*ast.IntegerLiteral)
	assert.True(t, ok, "il is not *ast.IntegerLiteral got=%T", il)

	assert.Equal(t, integ.Value, value)
	assert.Equal(t, integ.TokenLiteral(), fmt.Sprintf("%d", value))
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) {
	ident, ok := expr.(*ast.Identifier)
	assert.True(t, ok, "expr is not an *ast.Identifier")
	assert.Equal(t, value, ident.Value, "ident.Value not %s, got=%s", value, ident.Value)
	assert.Equal(t, value, ident.TokenLiteral(), "ident.TokenLiteral not %s, got=%s", value, ident.TokenLiteral())
}

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	assert.Equal(t, "let", s.TokenLiteral())

	stmt, ok := s.(*ast.LetStatement)
	assert.True(t, ok, "s not *ast.LetStatement, got=%T", s)

	assert.Equal(t, name, stmt.Name.Value, "stmt.Name.Value not '%s', got=%s", name, stmt.Name.Value)
	assert.Equal(t, name, stmt.Name.TokenLiteral(), "stmt.Name.TokenLiteral() not '%s'. got=%s", name, stmt.Name.TokenLiteral())
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	for _, msg := range errs {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
