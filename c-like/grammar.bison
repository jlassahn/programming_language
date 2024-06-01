
%{
#include <stdio.h>

int yylex(void);
void yyerror(const char *s);

%}

%token IDENTIFIER;
%token NUMBER;
%token CHARCONST;
%token STRINGCONST;

// keywords
%token AS
%token CONSTANT
%token IMPORT
%token LINKAGE
%token POINTER
%token PRIVATE
%token SIZEOF
%token TYPEDEF
%token USING

// operators
%token ASSIGNMENT_OP /* = *= /= %= += -= <<= >>= &= ^= |= */
%token LOG_OR_OP /* || */
%token LOG_AND_OP /* && */
%token OR_OP /* | */
%token AND_OP /* & */
%token XOR_OP /* ^ */
%token EQUAL_OP /* == != */
%token RELATIONAL_OP /* < > <= >= */
%token SHIFT_OP /* << >> */
%token ADD_OP /* + - */
%token MULT_OP /* * / % */
%token PREFIX_OP /* & * + - ~ ! */
%token INC_OP /* ++ -- */

%%

file:
  /* empty */
| file file_element
// FIXME maybe a file header?
;

file_element:
  import_statement
| using_statement
| external_declaration
;

import_statement:
  IMPORT namespace_expression ';'
| IMPORT PRIVATE namespace_expression ';'
;

using_statement:
  USING namespace_expression ';'
| USING namespace_expression AS namespace_expression ';'
;

external_declaration:
  declaration_type IDENTIFIER variable_properties initializer ';'
| declaration_type IDENTIFIER '(' proto_params ')' function_properties ';'
| declaration_type IDENTIFIER '(' proto_params ')' function_properties compound_statement
// FIXME struct ...
// FIXME union ...
// FIXME enum ...
;

proto_params:
  /* empty */
| proto_param_list
| proto_param_list ',' // extra trailing comma allowed
;

proto_param_list:
  proto_param
| proto_param_list ',' proto_param
;

proto_param:
  type_expression
| type_expression IDENTIFIER initializer
;

compound_statement:
  '{' '}'  // FIXME needs contents
;

declaration_type:
  type_expression
| storage_specifier type_expression
;

initializer:
  /* empty */
| '=' expression
| '=' '{' '}' // FIXME needs contents
;

expression:
  assignment_expression
// no comma operator
;

assignment_expression:
  conditional_expression
| unary_expression ASSIGNMENT_OP assignment_expression
	// probably doesn't have to be unary, could be 
	// conditional_expression ASSIGNMENT_OP assignment_expression
;

conditional_expression:
  logical_or_expression
| logical_or_expression '?' expression ':' conditional_expression
;

logical_or_expression:
  logical_and_expression
| logical_or_expression LOG_OR_OP logical_and_expression
;

logical_and_expression:
  inclusive_or_expression
| logical_and_expression LOG_AND_OP inclusive_or_expression
;

inclusive_or_expression:
  exclusive_or_expression
| inclusive_or_expression OR_OP exclusive_or_expression
;

exclusive_or_expression:
  and_expression
| exclusive_or_expression XOR_OP and_expression
;

and_expression:
  equality_expression
| and_expression AND_OP equality_expression
;

equality_expression:
  relational_expression
| equality_expression EQUAL_OP relational_expression
;

relational_expression:
  shift_expression
| relational_expression RELATIONAL_OP shift_expression
;

shift_expression:
  additive_expression
| shift_expression SHIFT_OP additive_expression
;

additive_expression:
  multiplicative_expression
| additive_expression ADD_OP multiplicative_expression
;

multiplicative_expression:
  unary_expression
| multiplicative_expression MULT_OP unary_expression
  //no cast expression, casts look like function calls now
;

unary_expression:
  postfix_expression
| SIZEOF unary_expression
| PREFIX_OP unary_expression
| INC_OP unary_expression
;

postfix_expression:
  primary_expression
| postfix_expression '[' ']' // FIXME contents
| postfix_expression '(' ')' // FIXME contents
| postfix_expression '{' '}' // FIXME contents
| postfix_expression '.' IDENTIFIER // ambiguous with namespace_expression
| postfix_expression INC_OP
// no arrow operator postfix_expression -> IDENTIFIER
;

primary_expression:
  type_expression
| NUMBER
| CHARCONST
| STRINGCONST
| '(' expression ')'
;

type_expression:
  namespace_expression
| type_modifier type_expression
;

	// ambiguity resolution:
	// Conflict is between the namespace selection dot and
	// the struct member dot.
	// They can be treated as equivalent in non-type expressions.
	// For types, the namespace dot needs to have higher precedence than
	// the TYPE_MODIFIER operators, but a type can never have a 
	// struct member, so for types all dots are namespace dots.
	// So it works to do a greedy parse whenre the namespace expression
	// always claims as many dots as possible.
	// (as long as the later compiler stages don't care about the
	// difference between a dot operator in a NAMESPACE_EXPRESSION node
	// and the dot operator in a POSTFIX_EXPRESSION node)

namespace_expression:
  IDENTIFIER
| namespace_expression '.' IDENTIFIER // ambiguous with postfix_expression
;

type_modifier:
	POINTER
//FIXME more...
;

storage_specifier:
  CONSTANT
| TYPEDEF
//FIXME more...
;

variable_properties:
 /* empty */
| LINKAGE
// FIXME more...
;

function_properties:
 /* empty */
| LINKAGE
// FIXME more...
;

%%

int yylex(void)
{
	return 0;
}

void yyerror(const char *s)
{
	printf("ERROR: %s\n", s);
}

int main(void)
{
	printf("return value = %d\n", yyparse());
	return 0;
}

