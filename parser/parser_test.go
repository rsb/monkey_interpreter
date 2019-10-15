package parser_test

import (
	"testing"

	"github.com/rsb/monkey_interpreter/ast"
	"github.com/rsb/monkey_interpreter/lexer"
	"github.com/rsb/monkey_interpreter/parser"
	"github.com/rsb/monkey_interpreter/token"

	"github.com/stretchr/testify/assert"
)

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

func TestLetStatements(t *testing.T) {
	assert := assert.New(t)
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	assert.NotNil(program)
	assert.Equal(3, len(program.Statements))

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		assert.Equal(stmt.TokenLiteral(), "let")
		letStmt, ok := stmt.(*ast.LetStatement)
		assert.True(ok)
		assert.Equal(letStmt.Name.Value, tt.expectedIdentifier)
		assert.Equal(letStmt.Name.TokenLiteral(), tt.expectedIdentifier)
	}
}

func TestLetSTatementsWithErrors(t *testing.T) {
	assert := assert.New(t)
	input := `
		let x 5;
		let = 10;
		let 838383;
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	assert.NotNil(program)

	errList := p.Errors()
	assert.Equal(len(errList), 3)
	assert.Equal(errList[0], "expected next token to be =, got INT instead")
	assert.Equal(errList[1], "expected next token to be IDENT, got = instead")
	assert.Equal(errList[2], "expected next token to be IDENT, got INT instead")
}

func TestReturnStatements(t *testing.T) {
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
