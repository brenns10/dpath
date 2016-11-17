package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const eof = 0

func TestNumeric(t *testing.T) {
	var sym yySymType
	uut := "1 123 1.2 .2 12.34 1.23e4 .23e4"
	l := NewLexer(strings.NewReader(uut))
	assert.Equal(t, l.Lex(&sym), INTEGER_LITERAL)
	assert.Equal(t, l.Lex(&sym), INTEGER_LITERAL)
	assert.Equal(t, l.Lex(&sym), DECIMAL_LITERAL)
	assert.Equal(t, l.Lex(&sym), DECIMAL_LITERAL)
	assert.Equal(t, l.Lex(&sym), DECIMAL_LITERAL)
	assert.Equal(t, l.Lex(&sym), DOUBLE_LITERAL)
	assert.Equal(t, l.Lex(&sym), DOUBLE_LITERAL)
	assert.Equal(t, l.Lex(&sym), eof)
}

func TestString(t *testing.T) {
	var sym yySymType
	uut := "'abc' 'abc''def' \"abc\" \"abc\"\"def\""
	l := NewLexer(strings.NewReader(uut))
	assert.Equal(t, l.Lex(&sym), STRING_LITERAL)
	assert.Equal(t, l.Lex(&sym), STRING_LITERAL)
	assert.Equal(t, l.Lex(&sym), STRING_LITERAL)
	assert.Equal(t, l.Lex(&sym), STRING_LITERAL)
	assert.Equal(t, l.Lex(&sym), eof)
}

func TestKeywords(t *testing.T) {
	var sym yySymType
	uut := "or and idiv div mod eq ne lt le gt ge file dir to"
	l := NewLexer(strings.NewReader(uut))
	assert.Equal(t, l.Lex(&sym), OR)
	assert.Equal(t, l.Lex(&sym), AND)
	assert.Equal(t, l.Lex(&sym), INTEGER_DIVIDE)
	assert.Equal(t, l.Lex(&sym), DIVIDE)
	assert.Equal(t, l.Lex(&sym), MODULUS)
	assert.Equal(t, l.Lex(&sym), VEQ)
	assert.Equal(t, l.Lex(&sym), VNE)
	assert.Equal(t, l.Lex(&sym), VLT)
	assert.Equal(t, l.Lex(&sym), VLE)
	assert.Equal(t, l.Lex(&sym), VGT)
	assert.Equal(t, l.Lex(&sym), VGE)
	assert.Equal(t, l.Lex(&sym), FILE)
	assert.Equal(t, l.Lex(&sym), DIR)
	assert.Equal(t, l.Lex(&sym), TO)
	assert.Equal(t, l.Lex(&sym), eof)
}

func TestSymbols(t *testing.T) {
	var sym yySymType
	uut := ":: $ ( ) [ ] , + - * / = != < <= > >= @ .. ."
	l := NewLexer(strings.NewReader(uut))
	assert.Equal(t, l.Lex(&sym), AXIS)
	assert.Equal(t, l.Lex(&sym), DOLLAR)
	assert.Equal(t, l.Lex(&sym), LPAREN)
	assert.Equal(t, l.Lex(&sym), RPAREN)
	assert.Equal(t, l.Lex(&sym), LBRACKET)
	assert.Equal(t, l.Lex(&sym), RBRACKET)
	assert.Equal(t, l.Lex(&sym), COMMA)
	assert.Equal(t, l.Lex(&sym), PLUS)
	assert.Equal(t, l.Lex(&sym), MINUS)
	assert.Equal(t, l.Lex(&sym), MULTIPLY)
	assert.Equal(t, l.Lex(&sym), SLASH)
	assert.Equal(t, l.Lex(&sym), GEQ)
	assert.Equal(t, l.Lex(&sym), GNE)
	assert.Equal(t, l.Lex(&sym), GLT)
	assert.Equal(t, l.Lex(&sym), GLE)
	assert.Equal(t, l.Lex(&sym), GGT)
	assert.Equal(t, l.Lex(&sym), GGE)
	assert.Equal(t, l.Lex(&sym), ATTR)
	assert.Equal(t, l.Lex(&sym), DOTDOT)
	assert.Equal(t, l.Lex(&sym), DOT)
	assert.Equal(t, l.Lex(&sym), eof)
}

func TestQname(t *testing.T) {
	var sym yySymType
	uut := "legal legal123 More_Legal_Names Me-.too _also"
	l := NewLexer(strings.NewReader(uut))
	assert.Equal(t, l.Lex(&sym), QNAME)
	assert.Equal(t, l.Lex(&sym), QNAME)
	assert.Equal(t, l.Lex(&sym), QNAME)
	assert.Equal(t, l.Lex(&sym), QNAME)
	assert.Equal(t, l.Lex(&sym), QNAME)
	assert.Equal(t, l.Lex(&sym), eof)
}

func TestInvalidNames(t *testing.T) {
	assertDoesNotLex(t, "0abc", QNAME)
	assertDoesNotLex(t, "-abc", QNAME)
}
