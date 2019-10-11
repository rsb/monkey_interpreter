package parser_test

import (
	"testing"

	"github.com/rsb/monkey_interpreter/ast"
	"github.com/rsb/monkey_interpreter/lexer"
	"github.com/rsb/monkey_interpreter/parser"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(len(program.Statements), 3)

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
