%{
package parse

import (
	"framecc/ast"
)

%}

%union {
    nsl []ast.Node
    n ast.Node
    ival int
    sval string
    svalarr []string
}

%token '{' '}' '[' ']' '<' '>' ','

%token tWhitespace

%token <sval> tIdentifier
%token tNumber

%token tStruct tType
%token tStringType
%token <n> tIntegerType
%token <sval> tIntegerFlag

%token tEOL

%type <n> type int_type
%type <n> field
%type <n> field_list
%type <n> scope
%type <n> reference
%type <n> struct anon_struct decl
%type <n> decl_list
%type <svalarr> int_modifiers

%start file

%%

file
    : ws✳ decl_list ws✳
      { yylex.(*Lexer).Ast().Scope = $2.(*ast.Scope) }
    ;

decl_list
    : ws✳ { $$ = ast.NewScope() }
    | decl_list ws✳ decl ws✳
      { $1.(*ast.Scope).Add($3.(*ast.Struct)) }
    ;

decl
	: tType struct
      {
          $$ = $2
              yylex.(*Lexer).AddDecl($$)
      }

struct
	: tIdentifier anon_struct
      {
          $2.(*ast.Struct).Name = $1
          $$ = $2
      }
    ;

anon_struct
    : tStruct scope
      {
          $$ = &ast.Struct{
              Name: "AnonStruct_X",
              Scope: $2.(*ast.Scope),
          }
      }
    ;

scope
    : '{' field_list '}'
      { $$ = $2 }
    ;

field_list
    : ws✳ { $$ = ast.NewScope() }
    | field_list field
      { $1.(*ast.Scope).Add($2.(ast.Node)) }
    ;

field
    : tIdentifier type
      {
          $$ = &ast.Field{
              Name: $1,
              Type: $2.(ast.Node),
          }
      }
    ;

int_type
    : tIntegerType
    | tIntegerType '<' int_modifiers '>'
      {
          $1.(*ast.IntegerType).Modifiers = $3
          $$ = $1
      }
    ;

int_modifiers
	: ws✳ { $$ = make([]string, 0) }
    | tIntegerFlag
      { $$ = append($$, $1) }
    | int_modifiers ',' tIntegerFlag
      { $$ = append($$, $3) }
    ;

reference
	: tIdentifier
      {
	      $$ = &ast.DeclReference{
			  DeclName: $1,
	      }
      }
    ;

type
    : int_type
    | anon_struct
    | reference
    ;



ws
	: tWhitespace
	;

ws＋
	: ws
	| ws ws＋

ws✳
	: /* empty */
	| ws＋

%%
