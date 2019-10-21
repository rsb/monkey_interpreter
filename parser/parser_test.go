package parser_test

import (
	"testing"

	"github.com/rsb/monkey_interpreter/ast"
	"github.com/rsb/monkey_interpreter/lexer"
	"github.com/rsb/monkey_interpreter/parser"
	"github.com/rsb/monkey_interpreter/token"

	"github.com/stretchr/testify/assert"
)

func TestLetStatements(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		input                string
		expectedIdentitifier string
		expectedValue        interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10;", "y", 10},
		{"let foobar = 83383;", "foobar", 83383},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		assert.NotNil(program)
		assert.Len(program.Statements, 1, "program should only have 1 statement")

		stmt := program.Statements[0]
		testLetStatement(t, stmt, tt.expectedIdentitifier)

		val := stmt.(*ast.LetStatement).Value
		testLiteralExpression(t, val, tt.expectedValue)
	}
}

func TestLetStatementErrorNoEqual(t *testing.T) {
	assert := assert.New(t)

	input := `let x 5`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	assert.NotNil(program)
	errList := p.Errors()
	assert.Equal(1, len(errList))
	assert.Equal(errList[0], "expected next token to be =, got INT instead")
}

func TestLetStatementErrorNoIDENTAndNoEqual(t *testing.T) {
	assert := assert.New(t)

	input := `let 838383`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	assert.NotNil(program)
	errList := p.Errors()
	assert.Equal(1, len(errList))
	assert.Equal(errList[0], "expected next token to be IDENT, got INT instead")
}

func TestLetStatementErrorNoIDENT(t *testing.T) {
	assert := assert.New(t)

	input := `let = 10`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	assert.NotNil(program)
	errList := p.Errors()

	assert.Equal(2, len(errList))
	assert.Equal("expected next token to be IDENT, got = instead", errList[0])
	assert.Equal("no prefix parse function for = found", errList[1])
}
func TestThreeReturnStatements(t *testing.T) {
	assert := assert.New(t)
	input := `
	return 5;
	return 10;
	return 9834443;
	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	assert.NotNil(program)
	assert.Equal(3, len(program.Statements))
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		assert.True(ok)
		assert.Equal(returnStmt.TokenLiteral(), "return")
	}
}

func TestNodeStringMethod(t *testing.T) {
	assert := assert.New(t)
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	assert.Equal("let myVar = anotherVar;", program.String())

}

func TestIdentifierExpression(t *testing.T) {
	assert := assert.New(t)
	in := "foobar;"
	l := lexer.New(in)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.Equal(1, len(program.Statements), "program has not enough statements")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatment)
	assert.True(ok)

	ident, ok := stmt.Expression.(*ast.Identifier)
	assert.True(ok)
	assert.Equal(ident.Value, "foobar")
	assert.Equal(ident.TokenLiteral(), "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	assert := assert.New(t)
	in := "5;"
	l := lexer.New(in)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.Equal(1, len(program.Statements), "program has not enough statements")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatment)
	assert.True(ok)

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	assert.True(ok)
	assert.Equal(ident.Value, int64(5))
	assert.Equal(ident.TokenLiteral(), "5")
}

func TestPrefixExpressions(t *testing.T) {
	assert := assert.New(t)
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		assert.Len(program.Statements, 1, "program should have 1 statement got %d", len(program.Statements))

		stmt, ok := program.Statements[0].(*ast.ExpressionStatment)
		assert.True(ok, "stmt is not *ast.ExpressionStatement got=%T", program.Statements[0])

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		assert.True(ok, "exp is not a *ast.PrefixEpression got=%T", stmt.Expression)

		assert.Equal(exp.Operator, tt.operator)
		testIntegerLiteral(t, exp.Right, tt.integerValue)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	assert := assert.New(t)

	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		assert.Len(program.Statements, 1, "program should have 1 statement got %d", len(program.Statements))

		stmt, ok := program.Statements[0].(*ast.ExpressionStatment)
		assert.True(ok, "stmt is not *ast.ExpressionStatement got=%T", program.Statements[0])

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestParseProgram_OperatorPrecendence(t *testing.T) {
	assert := assert.New(t)

	infixTests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		assert.Equal(tt.expected, actual)
	}
}
