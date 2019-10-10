package lexer_test

import (
	"testing"

	"github.com/rsb/monkey_interpreter/lexer"
	"github.com/rsb/monkey_interpreter/token"

	"github.com/stretchr/testify/assert"
)

func TestNextTokenRandomSimpleTokens(t *testing.T) {
	assert := assert.New(t)
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := lexer.New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		assert.Equal(tok.Type, tt.expectedType)
		assert.Equal(tok.Literal, tt.expectedLiteral)
	}
}

func TestNextTokenStatements(t *testing.T) {
	assert := assert.New(t)
	input := `let five = 5;
	let ten = 10;

	let add = fn(x,y) {
		x + y;
	};

	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := lexer.New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		assert.Equal(tok.Type, tt.expectedType)
		assert.Equal(tok.Literal, tt.expectedLiteral)
	}
}

func TestNextTokenStatementsWithIfBlocks(t *testing.T) {
	assert := assert.New(t)
	input := `if (5 < 10) {
		return true;
	} else {
		return false;
	}
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := lexer.New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		assert.Equal(tok.Type, tt.expectedType)
		assert.Equal(tok.Literal, tt.expectedLiteral)
	}
}

func TestNextTokenStatementsWithPeekAHead(t *testing.T) {
	assert := assert.New(t)
	input := `10 == 10;
		10 != 9;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := lexer.New(input)
	for _, tt := range tests {
		tok := l.NextToken()
		assert.Equal(tok.Type, tt.expectedType)
		assert.Equal(tok.Literal, tt.expectedLiteral)
	}
}
