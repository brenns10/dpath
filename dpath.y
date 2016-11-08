%{
package main
%}

%union {
    tree *ParseTree
}

%token STRING_LITERAL
%token INTEGER_LITERAL
%token DECIMAL_LITERAL
%token DOUBLE_LITERAL

%token OR
%token AND
%token DIVIDE
%token INTEGER_DIVIDE
%token MODULUS
%token VEQ
%token VNE
%token VLT
%token VLE
%token VGT
%token VGE
%token IS
%token UNION
%token FILE
%token DIR
%token TO
%token AXIS

%token QNAME

%token DOLLAR
%token LPAREN
%token RPAREN
%token LBRACKET
%token RBRACKET
%token COMMA
%token PLUS
%token MINUS
%token MULTIPLY
%token SLASH
%token GEQ
%token GNE
%token GLT
%token GLE
%token GGT
%token GGE
%token UNIONSYM
%token ATTR
%token DOTDOT
%token DOT

%%
XPath: Expr;

Expr: ExprSingle
      | Expr COMMA ExprSingle
      ;

ExprSingle: OrExpr
            ;

OrExpr: AndExpr
        | AndExpr OR AndExpr
        ;

AndExpr: ComparisonExpr
         | ComparisonExpr AND ComparisonExpr
         ;

ComparisonExpr: RangeExpr
                | RangeExpr ValueComp RangeExpr
                | RangeExpr GeneralComp RangeExpr
                | RangeExpr NodeComp RangeExpr
                ;

ValueComp: VEQ | VNE | VLT | VLE | VGT | VGE;

GeneralComp: GEQ | GNE | GLT | GLE | GGT | GGE;
NodeComp: IS;

RangeExpr: AdditiveExpr
           | AdditiveExpr TO AdditiveExpr
           ;

AdditiveExpr: MultiplicativeExpr
              | MultiplicativeExpr PLUS MultiplicativeExpr
              | MultiplicativeExpr MINUS MultiplicativeExpr
              ;

MultiplicativeExpr: UnaryExpr
                    | UnaryExpr MULTIPLY UnaryExpr
                    | UnaryExpr DIVIDE UnaryExpr
                    | UnaryExpr INTEGER_DIVIDE UnaryExpr
                    | UnaryExpr MODULUS UnaryExpr
                    ;

UnaryExpr: ValueExpr
           | PLUS ValueExpr
           | MINUS ValueExpr
           ;

ValueExpr: PathExpr
           ;

PathExpr: RelativePathExpr
          | SLASH RelativePathExpr
          | SLASH SLASH RelativePathExpr
          ;

RelativePathExpr: StepExpr
                  | RelativePathExpr SLASH StepExpr
                  | RelativePathExpr SLASH SLASH StepExpr
                  ;

StepExpr: AxisStep
          | FilterExpr
          ;

AxisStep: NodeStep
          | NodeStep PredicateList

NodeStep: QNAME AXIS NodeTest
          | ATTR NodeTest
          | DOTDOT
          | NodeTest
          ;

NodeTest: KindTest
          | NameTest
          ;

NameTest: QNAME
          | MULTIPLY
          ;

KindTest: FILE LPAREN RPAREN
          | DIR LPAREN RPAREN
          ;

PredicateList: Predicate
               | PredicateList Predicate
               ;

Predicate: LBRACKET Expr RBRACKET
           ;

FilterExpr: PrimaryExpr
            | PrimaryExpr PredicateList
            ;

PrimaryExpr: Literal
             | ParenthesizedExpr
             | ContextItemExpr
             | FunctionCall
             ;

ParenthesizedExpr: LPAREN Expr RPAREN
                   | LPAREN RPAREN
                   ;

ContextItemExpr: DOT
                 ;

FunctionCall: QNAME LPAREN RPAREN
              | QNAME LPAREN ArgumentList RPAREN
              ;

ArgumentList: ExprSingle
              | ArgumentList COMMA ExprSingle
              ;

Literal:  STRING_LITERAL
          | INTEGER_LITERAL
          | DECIMAL_LITERAL
          | DOUBLE_LITERAL
          ;

%%
