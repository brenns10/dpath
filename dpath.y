%{
package main
%}

%union {
    tree ParseTree
    str string
    num int
    args []ParseTree
}

%token  <str>           STRING_LITERAL
%token  <str>           INTEGER_LITERAL
%token  <str>           DECIMAL_LITERAL
%token  <str>           DOUBLE_LITERAL
%token  <str>           QNAME

%token  <num>           OR
%token  <num>           AND
%token  <num>           DIVIDE
%token  <num>           INTEGER_DIVIDE
%token  <num>           MODULUS
%token  <num>           VEQ
%token  <num>           VNE
%token  <num>           VLT
%token  <num>           VLE
%token  <num>           VGT
%token  <num>           VGE
%token  <num>           IS
%token  <num>           UNION
%token  <num>           FILE
%token  <num>           DIR
%token  <num>           TO
%token  <num>           AXIS

%token  <num>           DOLLAR
%token  <num>           LPAREN
%token  <num>           RPAREN
%token  <num>           LBRACKET
%token  <num>           RBRACKET
%token  <num>           COMMA
%token  <num>           PLUS
%token  <num>           MINUS
%token  <num>           MULTIPLY
%token  <num>           SLASH
%token  <num>           GEQ
%token  <num>           GNE
%token  <num>           GLT
%token  <num>           GLE
%token  <num>           GGT
%token  <num>           GGE
%token  <num>           UNIONSYM
%token  <num>           ATTR
%token  <num>           DOTDOT
%token  <num>           DOT

%type   <tree>          XPath
%type   <tree>          Expr
%type   <tree>          ExprSingle
%type   <tree>          OrExpr
%type   <tree>          AndExpr
%type   <tree>          ComparisonExpr
%type   <str>           ValueComp
%type   <str>           GeneralComp
%type   <str>           NodeComp
%type   <tree>          RangeExpr
%type   <tree>          AdditiveExpr
%type   <tree>          MultiplicativeExpr
%type   <tree>          UnaryExpr
%type   <tree>          ValueExpr
%type   <tree>          PathExpr
%type   <args>          RelativePathExpr
%type   <tree>          StepExpr
%type   <tree>          AxisStep
%type   <tree>          NodeStep
%type   <tree>          NodeTest
%type   <tree>          NameTest
%type   <tree>          KindTest
%type   <args>          PredicateList
%type   <tree>          Predicate
%type   <tree>          FilterExpr
%type   <tree>          PrimaryExpr
%type   <tree>          ParenthesizedExpr
%type   <tree>          ContextItemExpr
%type   <tree>          FunctionCall
%type   <args>          ArgumentList
%type   <tree>          Literal

%%
XPath:          Expr {parserResult = $1}
                ;

Expr:           ExprSingle {$$ = $1}
        |       Expr COMMA ExprSingle
                ;

ExprSingle:     OrExpr {$$ = $1}
                ;

OrExpr:         AndExpr {$$ = $1}
        |       AndExpr OR AndExpr {$$ = newBinopTree("or", $1, $3)}
                ;

AndExpr:        ComparisonExpr {$$ = $1}
        |       ComparisonExpr AND ComparisonExpr {$$ = newBinopTree("and", $1, $3)}
                ;

ComparisonExpr: RangeExpr {$$ = $1}
        |       RangeExpr ValueComp RangeExpr {$$ = newBinopTree($2, $1, $3)}
        |       RangeExpr GeneralComp RangeExpr {$$ = newBinopTree($2, $1, $3)}
        |       RangeExpr NodeComp RangeExpr {$$ = newBinopTree($2, $1, $3)}
                ;

ValueComp:      VEQ {$$ = "eq"}
        |       VNE {$$ = "ne"}
        |       VLT {$$ = "lt"}
        |       VLE {$$ = "le"}
        |       VGT {$$ = "gt"}
        |       VGE {$$ = "ge"}
                ;

GeneralComp:    GEQ {$$ = "="}
        |       GNE {$$ = "!="}
        |       GLT {$$ = "<"}
        |       GLE {$$ = "<="}
        |       GGT {$$ = ">"}
        |       GGE {$$ = ">="}
                ;

NodeComp:       IS {$$ = "is"}
                ;

RangeExpr:      AdditiveExpr {$$ = $1}
        |       AdditiveExpr TO AdditiveExpr {$$ = newBinopTree("to", $1, $3)}
                ;

AdditiveExpr:   MultiplicativeExpr {$$ = $1}
        |       MultiplicativeExpr PLUS MultiplicativeExpr {$$ = newBinopTree("+", $1, $3)}
        |       MultiplicativeExpr MINUS MultiplicativeExpr {$$ = newBinopTree("-", $1, $3)}
                ;

MultiplicativeExpr:
                UnaryExpr {$$ = $1}
        |       UnaryExpr MULTIPLY UnaryExpr {$$ = newBinopTree("*", $1, $3)}
        |       UnaryExpr DIVIDE UnaryExpr {$$ = newBinopTree("div", $1, $3)}
        |       UnaryExpr INTEGER_DIVIDE UnaryExpr {$$ = newBinopTree("idiv", $1, $3)}
        |       UnaryExpr MODULUS UnaryExpr {$$ = newBinopTree("mod", $1, $3)}
                ;

UnaryExpr:      ValueExpr {$$ = $1}
        |       PLUS ValueExpr {$$ = newUnopTree("+", $2)}
        |       MINUS ValueExpr {$$ = newUnopTree("-", $2)}
                ;

ValueExpr:      PathExpr {$$ = $1}
                ;

PathExpr:       RelativePathExpr
                {
                    if len($1) == 1 {
                        $$ = $1[0]
                    } else {
                        $$ = newPathTree($1, false)
                    }
                }
        |       SLASH RelativePathExpr {$$ = newPathTree($2, true)}
        |       SLASH SLASH RelativePathExpr {$$ = newPathTree(append([]ParseTree{nil}, $3...), true)}
                ;

RelativePathExpr:
                StepExpr {$$ = []ParseTree{$1}}
        |       RelativePathExpr SLASH StepExpr {$$ = append($1, $3)}
        |       RelativePathExpr SLASH SLASH StepExpr {$$ = append($1, nil, $4)}
                ;

StepExpr:       AxisStep {$$ = $1}
        |       FilterExpr {$$ = $1}
                ;

AxisStep:       NodeStep {$$ = $1}
        |       NodeStep PredicateList {$$ = newFilteredSequenceTree($1, $2)}
                ;

NodeStep:       QNAME AXIS NodeTest {$$ = newAxisTree($1, $3)}
        |       ATTR NodeTest {$$ = newAxisTree("attr", $2)}
        |       DOTDOT {$$ = newKindTree("..")}
        |       NodeTest {$$ = $1}
                ;

NodeTest:       KindTest {$$ = $1}
        |       NameTest {$$ = $1}
                ;

NameTest:       QNAME {$$ = newNameTree($1)}
        |       MULTIPLY {$$ = newKindTree("*")}
                ;

KindTest:       FILE LPAREN RPAREN {$$ = newKindTree("file")}
        |       DIR LPAREN RPAREN {$$ = newKindTree("dir")}
                ;

PredicateList:  Predicate {$$ = []ParseTree{$1}}
        |       PredicateList Predicate {$$ = append($1, $2)}
                ;

Predicate:      LBRACKET Expr RBRACKET {$$ = $2}
                ;

FilterExpr:     PrimaryExpr {$$ = $1}
        |       PrimaryExpr PredicateList {$$ = newFilteredSequenceTree($1, $2)}
                ;

PrimaryExpr:    Literal {$$ = $1}
        |       ParenthesizedExpr {$$ = $1}
        |       ContextItemExpr {$$ = $1}
        |       FunctionCall {$$ = $1}
                ;

ParenthesizedExpr:
                LPAREN Expr RPAREN {$$ = $2}
        |       LPAREN RPAREN {$$ = newEmptySequenceTree()}
                ;

ContextItemExpr:DOT {$$ = newContextItemTree()}
                ;

FunctionCall:   QNAME LPAREN RPAREN {$$ = newFunccallTree($1, []ParseTree{})}
        |       QNAME LPAREN ArgumentList RPAREN {$$ = newFunccallTree($1, $3)}
                ;

ArgumentList:   ExprSingle {$$ = []ParseTree{$1}}
        |       ArgumentList COMMA ExprSingle {$$ = append($1, $3)}
                ;

Literal:        STRING_LITERAL {$$ = newStringTree($1)}
        |       INTEGER_LITERAL {$$ = newIntegerTree($1)}
        |       DECIMAL_LITERAL {$$ = newDoubleTree($1)}
        |       DOUBLE_LITERAL {$$ = newDoubleTree($1)}
                ;

%%
