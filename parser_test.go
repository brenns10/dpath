package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertParses(t *testing.T, s string) {
	_, e := ParseString(s)
	assert.Nil(t, e)
}

func TestSimpleExpressionsParse(t *testing.T) {
	assertParses(t, "1")
	assertParses(t, "1.23")
	assertParses(t, "1.0E-1")
	assertParses(t, "'yo'")
	assertParses(t, "identifier")
	assertParses(t, "()")
	assertParses(t, "1, 2")
}

func TestRangeExpressions(t *testing.T) {
	assertParses(t, "1 + 1 to 2 + 2")
	assertParses(t, "//* to //*")
}

func TestLogicExpressions(t *testing.T) {
	assertParses(t, "x = 5 and 3 + 3 eq 6")
	assertParses(t, "10 or .")
}

func TestComparisons(t *testing.T) {
	assertParses(t, "1+1 eq 1")
	assertParses(t, "1 = 1")
	assertParses(t, "1 ne 1* /*")
	assertParses(t, "1 != 1")
	assertParses(t, "1 gt 1")
	assertParses(t, ". > /name[@blah]")
	assertParses(t, "1 ge 1")
	assertParses(t, "1 >= 1")
	assertParses(t, "1 lt 1")
	assertParses(t, "1 < 1")
	assertParses(t, "1 le 1")
	assertParses(t, "-1 <= 1")
	assertParses(t, "/Blah is .")
}

func TestStepExprs(t *testing.T) {
	assertParses(t, "/file()")
	assertParses(t, "//dir()")
	assertParses(t, "file()[@owner = 'stephen']")
}
