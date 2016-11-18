//line dpath.y:2
package main

import __yyfmt__ "fmt"

//line dpath.y:2
//line dpath.y:5
type yySymType struct {
	yys  int
	tree ParseTree
	str  string
	num  int
	args []ParseTree
}

const STRING_LITERAL = 57346
const INTEGER_LITERAL = 57347
const DECIMAL_LITERAL = 57348
const DOUBLE_LITERAL = 57349
const QNAME = 57350
const OR = 57351
const AND = 57352
const DIVIDE = 57353
const INTEGER_DIVIDE = 57354
const MODULUS = 57355
const VEQ = 57356
const VNE = 57357
const VLT = 57358
const VLE = 57359
const VGT = 57360
const VGE = 57361
const FILE = 57362
const DIR = 57363
const TO = 57364
const AXIS = 57365
const DOLLAR = 57366
const POUND = 57367
const LPAREN = 57368
const RPAREN = 57369
const LBRACKET = 57370
const RBRACKET = 57371
const COMMA = 57372
const PLUS = 57373
const MINUS = 57374
const MULTIPLY = 57375
const SLASH = 57376
const GEQ = 57377
const GNE = 57378
const GLT = 57379
const GLE = 57380
const GGT = 57381
const GGE = 57382
const ATTR = 57383
const DOTDOT = 57384
const DOT = 57385

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"STRING_LITERAL",
	"INTEGER_LITERAL",
	"DECIMAL_LITERAL",
	"DOUBLE_LITERAL",
	"QNAME",
	"OR",
	"AND",
	"DIVIDE",
	"INTEGER_DIVIDE",
	"MODULUS",
	"VEQ",
	"VNE",
	"VLT",
	"VLE",
	"VGT",
	"VGE",
	"FILE",
	"DIR",
	"TO",
	"AXIS",
	"DOLLAR",
	"POUND",
	"LPAREN",
	"RPAREN",
	"LBRACKET",
	"RBRACKET",
	"COMMA",
	"PLUS",
	"MINUS",
	"MULTIPLY",
	"SLASH",
	"GEQ",
	"GNE",
	"GLT",
	"GLE",
	"GGT",
	"GGE",
	"ATTR",
	"DOTDOT",
	"DOT",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line dpath.y:235

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 79
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 258

var yyAct = [...]int{

	3, 17, 25, 15, 72, 9, 7, 2, 8, 6,
	68, 32, 33, 34, 35, 22, 5, 60, 61, 10,
	69, 109, 42, 63, 64, 65, 77, 38, 39, 42,
	73, 107, 41, 36, 102, 106, 82, 81, 12, 13,
	40, 16, 59, 84, 79, 62, 78, 44, 23, 24,
	37, 71, 87, 88, 86, 83, 43, 26, 38, 39,
	85, 110, 105, 41, 111, 42, 90, 91, 89, 75,
	96, 40, 76, 74, 98, 103, 99, 104, 101, 99,
	11, 100, 92, 93, 94, 95, 32, 33, 34, 35,
	22, 29, 28, 66, 67, 27, 21, 19, 30, 108,
	31, 20, 38, 39, 18, 14, 46, 41, 36, 80,
	45, 4, 112, 12, 13, 40, 16, 32, 33, 34,
	35, 22, 1, 23, 24, 37, 0, 0, 0, 0,
	0, 0, 0, 38, 39, 0, 0, 0, 41, 36,
	0, 0, 0, 0, 12, 13, 40, 16, 32, 33,
	34, 35, 22, 0, 23, 24, 37, 0, 0, 32,
	33, 34, 35, 22, 38, 39, 0, 0, 0, 41,
	36, 0, 0, 0, 0, 38, 39, 40, 97, 0,
	41, 36, 0, 0, 0, 23, 24, 37, 40, 70,
	32, 33, 34, 35, 22, 0, 23, 24, 37, 0,
	0, 32, 33, 34, 35, 22, 38, 39, 0, 0,
	0, 41, 36, 0, 0, 0, 0, 38, 39, 40,
	16, 0, 41, 36, 0, 0, 0, 23, 24, 37,
	40, 47, 48, 49, 50, 51, 52, 0, 23, 24,
	37, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 53, 54, 55, 56, 57, 58,
}
var yyPact = [...]int{

	113, -1000, -1, -1000, -1000, 47, 37, 217, 20, -14,
	12, -1000, 186, 186, -1000, -24, 155, -1000, -1000, -1000,
	2, 2, 46, 38, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 82, -1000, 11, 10,
	-1000, 51, 113, 113, 113, 113, 113, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 113,
	113, 113, 113, 113, 113, 113, -1000, -1000, 144, -24,
	197, 2, -1000, 113, 2, 38, 7, -1000, -1000, 35,
	-1000, 8, 4, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 197, -24, -1000,
	-8, -1000, -1000, 34, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 113, -1000,
}
var yyPgo = [...]int{

	0, 122, 7, 0, 111, 16, 9, 110, 106, 6,
	8, 5, 19, 80, 105, 3, 1, 104, 101, 2,
	100, 98, 51, 4, 97, 96, 95, 92, 91, 75,
	57,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 3, 4, 4, 5, 5, 6,
	6, 6, 7, 7, 7, 7, 7, 7, 8, 8,
	8, 8, 8, 8, 9, 9, 10, 10, 10, 11,
	11, 11, 11, 11, 12, 12, 12, 13, 14, 14,
	14, 15, 15, 15, 16, 16, 17, 17, 18, 18,
	18, 18, 19, 19, 20, 20, 20, 21, 21, 22,
	22, 23, 24, 24, 25, 25, 25, 25, 26, 26,
	27, 28, 28, 29, 29, 30, 30, 30, 30,
}
var yyR2 = [...]int{

	0, 1, 1, 3, 1, 1, 3, 1, 3, 1,
	3, 3, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 3, 1, 3, 3, 1,
	3, 3, 3, 3, 1, 2, 2, 1, 1, 2,
	3, 1, 3, 4, 1, 1, 1, 2, 3, 2,
	1, 1, 1, 1, 1, 1, 2, 3, 3, 1,
	2, 3, 1, 2, 1, 1, 1, 1, 3, 2,
	1, 3, 4, 1, 3, 1, 1, 1, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, -6, -9, -10, -11,
	-12, -13, 31, 32, -14, -15, 34, -16, -17, -24,
	-18, -25, 8, 41, 42, -19, -30, -26, -27, -28,
	-21, -20, 4, 5, 6, 7, 26, 43, 20, 21,
	33, 25, 30, 9, 10, -7, -8, 14, 15, 16,
	17, 18, 19, 35, 36, 37, 38, 39, 40, 22,
	31, 32, 33, 11, 12, 13, -13, -13, 34, -15,
	34, -22, -23, 28, -22, 23, 26, -19, 8, -2,
	27, 26, 26, 4, -3, -5, -6, -9, -9, -10,
	-11, -11, -12, -12, -12, -12, -16, 34, -15, -23,
	-2, -19, 27, -29, -3, 27, 27, 27, -16, 29,
	27, 30, -3,
}
var yyDef = [...]int{

	0, -2, 1, 2, 4, 5, 7, 9, 24, 26,
	29, 34, 0, 0, 37, 38, 0, 41, 44, 45,
	46, 62, 54, 0, 50, 51, 64, 65, 66, 67,
	52, 53, 75, 76, 77, 78, 0, 70, 0, 0,
	55, 0, 0, 0, 0, 0, 0, 12, 13, 14,
	15, 16, 17, 18, 19, 20, 21, 22, 23, 0,
	0, 0, 0, 0, 0, 0, 35, 36, 0, 39,
	0, 47, 59, 0, 63, 0, 0, 49, 54, 0,
	69, 0, 0, 56, 3, 6, 8, 10, 11, 25,
	27, 28, 30, 31, 32, 33, 42, 0, 40, 60,
	0, 48, 71, 0, 73, 68, 57, 58, 43, 61,
	72, 0, 74,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:87
		{
			parserResult = newSequenceTree(yyDollar[1].args)
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:90
		{
			yyVAL.args = []ParseTree{yyDollar[1].tree}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:91
		{
			yyVAL.args = append(yyDollar[1].args, yyDollar[3].tree)
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:94
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:97
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:98
		{
			yyVAL.tree = newBinopTree("or", yyDollar[1].tree, yyDollar[3].tree)
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:101
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:102
		{
			yyVAL.tree = newBinopTree("and", yyDollar[1].tree, yyDollar[3].tree)
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:105
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:106
		{
			yyVAL.tree = newBinopTree(yyDollar[2].str, yyDollar[1].tree, yyDollar[3].tree)
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:107
		{
			yyVAL.tree = newBinopTree(yyDollar[2].str, yyDollar[1].tree, yyDollar[3].tree)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:110
		{
			yyVAL.str = "eq"
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:111
		{
			yyVAL.str = "ne"
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:112
		{
			yyVAL.str = "lt"
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:113
		{
			yyVAL.str = "le"
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:114
		{
			yyVAL.str = "gt"
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:115
		{
			yyVAL.str = "ge"
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:118
		{
			yyVAL.str = "="
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:119
		{
			yyVAL.str = "!="
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:120
		{
			yyVAL.str = "<"
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:121
		{
			yyVAL.str = "<="
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:122
		{
			yyVAL.str = ">"
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:123
		{
			yyVAL.str = ">="
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:126
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:127
		{
			yyVAL.tree = newBinopTree("to", yyDollar[1].tree, yyDollar[3].tree)
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:130
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:131
		{
			yyVAL.tree = newBinopTree("+", yyDollar[1].tree, yyDollar[3].tree)
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:132
		{
			yyVAL.tree = newBinopTree("-", yyDollar[1].tree, yyDollar[3].tree)
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:136
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:137
		{
			yyVAL.tree = newBinopTree("*", yyDollar[1].tree, yyDollar[3].tree)
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:138
		{
			yyVAL.tree = newBinopTree("div", yyDollar[1].tree, yyDollar[3].tree)
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:139
		{
			yyVAL.tree = newBinopTree("idiv", yyDollar[1].tree, yyDollar[3].tree)
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:140
		{
			yyVAL.tree = newBinopTree("mod", yyDollar[1].tree, yyDollar[3].tree)
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:143
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line dpath.y:144
		{
			yyVAL.tree = newUnopTree("+", yyDollar[2].tree)
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line dpath.y:145
		{
			yyVAL.tree = newUnopTree("-", yyDollar[2].tree)
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:148
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:152
		{
			if len(yyDollar[1].args) == 1 {
				yyVAL.tree = yyDollar[1].args[0]
			} else {
				yyVAL.tree = newPathTree(yyDollar[1].args, false)
			}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line dpath.y:159
		{
			yyVAL.tree = newPathTree(yyDollar[2].args, true)
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:160
		{
			yyVAL.tree = newPathTree(append([]ParseTree{nil}, yyDollar[3].args...), true)
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:164
		{
			yyVAL.args = []ParseTree{yyDollar[1].tree}
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:165
		{
			yyVAL.args = append(yyDollar[1].args, yyDollar[3].tree)
		}
	case 43:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line dpath.y:166
		{
			yyVAL.args = append(yyDollar[1].args, nil, yyDollar[4].tree)
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:169
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:170
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:173
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line dpath.y:174
		{
			yyVAL.tree = newFilteredSequenceTree(yyDollar[1].tree, yyDollar[2].args)
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:177
		{
			yyVAL.tree = newAxisTree(yyDollar[1].str, yyDollar[3].tree)
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line dpath.y:178
		{
			yyVAL.tree = newAxisTree("attribute", yyDollar[2].tree)
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:179
		{
			yyVAL.tree = newKindTree("..")
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:180
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:183
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:184
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:187
		{
			yyVAL.tree = newNameTree(yyDollar[1].str)
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:188
		{
			yyVAL.tree = newKindTree("*")
		}
	case 56:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line dpath.y:189
		{
			yyVAL.tree = newNameTree(parseStringLiteral(yyDollar[2].str))
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:192
		{
			yyVAL.tree = newKindTree("file")
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:193
		{
			yyVAL.tree = newKindTree("dir")
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:196
		{
			yyVAL.args = []ParseTree{yyDollar[1].tree}
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line dpath.y:197
		{
			yyVAL.args = append(yyDollar[1].args, yyDollar[2].tree)
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:200
		{
			yyVAL.tree = newSequenceTree(yyDollar[2].args)
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:203
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 63:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line dpath.y:204
		{
			yyVAL.tree = newFilteredSequenceTree(yyDollar[1].tree, yyDollar[2].args)
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:207
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:208
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:209
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:210
		{
			yyVAL.tree = yyDollar[1].tree
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:214
		{
			yyVAL.tree = newSequenceTree(yyDollar[2].args)
		}
	case 69:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line dpath.y:215
		{
			yyVAL.tree = newEmptySequenceTree()
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:218
		{
			yyVAL.tree = newContextItemTree()
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:221
		{
			yyVAL.tree = newFunccallTree(yyDollar[1].str, []ParseTree{})
		}
	case 72:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line dpath.y:222
		{
			yyVAL.tree = newFunccallTree(yyDollar[1].str, yyDollar[3].args)
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:225
		{
			yyVAL.args = []ParseTree{yyDollar[1].tree}
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line dpath.y:226
		{
			yyVAL.args = append(yyDollar[1].args, yyDollar[3].tree)
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:229
		{
			yyVAL.tree = newStringTree(yyDollar[1].str)
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:230
		{
			yyVAL.tree = newIntegerTree(yyDollar[1].str)
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:231
		{
			yyVAL.tree = newDoubleTree(yyDollar[1].str)
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line dpath.y:232
		{
			yyVAL.tree = newDoubleTree(yyDollar[1].str)
		}
	}
	goto yystack /* stack new state and value */
}
