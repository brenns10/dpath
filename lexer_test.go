package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNumeric(t *testing.T) {
	var sym yySymType
	l := NewLexer(strings.NewReader("1 123 1.2 .2 12.34 1.23e4 .23e4"))
	assert.Equal(t, l.Lex(&sym), INTEGER_LITERAL)
	assert.Equal(t, l.Lex(&sym), INTEGER_LITERAL)
	assert.Equal(t, l.Lex(&sym), DECIMAL_LITERAL)
	assert.Equal(t, l.Lex(&sym), DECIMAL_LITERAL)
	assert.Equal(t, l.Lex(&sym), DECIMAL_LITERAL)
	assert.Equal(t, l.Lex(&sym), DOUBLE_LITERAL)
	assert.Equal(t, l.Lex(&sym), DOUBLE_LITERAL)
}
