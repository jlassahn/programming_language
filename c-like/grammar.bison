
%{
#include <stdio.h>

int yylex(void);
void yyerror(const char *s);

#define YYSTYPE double
#define YYCALLBACK MakeNode


double MakeNode(int len, double *args);

/*
// per rule behavior appears inside
// switch (yyn)
// in the output file
// byacc
#define MkNode yyval = MakeNode(yym, &yystack.l_mark[1-yym])
//bison
//#define MkNode yyval = MakeNode(yylen, &yyvsp[1-yylen])
*/

%}

/* punctuators
  ; = { } , . [ ] ( ) :
   conditional op
  ? :
   other uses of operator characters
  * for autodetect array size
*/

%token IDENTIFIER;
%token NUMBER;
%token CHARCONST;
%token STRINGCONST;

// keywords
%token ALIAS
%token ALLIGNMENT // FIXME maybe not real
%token AS
%token ARRAY
%token AUTO
%token BITFIELD
%token BREAK
%token CASE
%token CONSTANT
%token CONTINUE
%token DEFAULT
%token DO
%token ELSE
%token ENUM
%token FOR
%token GOTO
%token IF
%token IMPORT
%token INLINE
%token LINKAGE
%token LINKNAME
%token POINTER
%token PRIVATE
%token READONLY
%token REGISTER
%token RESTRICT
%token RETURN
%token SIZEOF
%token STATIC
%token STRUCT
%token SWITCH
%token TYPEDEF
%token UNION
%token USING
%token VOLATILE
%token WHILE
// reserve TEMPLATE and CLASS for future...

// operators
%token ASSIGNMENT_OP /* = *= /= %= += -= <<= >>= &= ^= |= */
%token LOG_OR_OP /* || */
%token LOG_AND_OP /* && */
%token OR_OP /* | */
%token AND_ADDR_OP /* & */
%token XOR_OP /* ^ */
%token EQUAL_OP /* == != */
%token RELATIONAL_OP /* < > <= >= */
%token SHIFT_OP /* << >> */
%token ADD_OP /* + - */
%token DIV_OP /* / % */
%token MULT_PTR_OP /* * */
%token NOT_OP /* ~ ! */
%token INC_OP /* ++ -- */

%token ELIPSIS

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
| IMPORT PRIVATE namespace_expression ';' { MkNode; }
;

using_statement:
  USING namespace_expression ';'
| USING namespace_expression AS namespace_expression ';'
;

proto_params:
  /* empty */
| proto_param_list
| proto_param_list ',' // extra trailing comma allowed
| proto_param_list ',' ELIPSIS
;

proto_param_list:
  proto_param
| proto_param_list ',' proto_param
;

proto_param:
  param_type
| param_type IDENTIFIER initializer
;

param_type:
  type_expression
| parameter_specifier type_expression
;

declaration_type:
  type_expression
| storage_specifier type_expression
;

initializer:
  /* empty */
| '=' expression
| '=' '{' initializers '}'
;

initializers:
  /* empty */
| initializer_list
| initializer_list ','  // extra trailing comma allowed
;

initializer_list:
  initializer_element
| initializer_list ',' initializer_element
;

initializer_element:
  '{' initializers '}'
| struct_initializer
| array_initializer
| constant_expression
;

struct_initializer:
  '.' IDENTIFIER '=' constant_expression
;

array_initializer:
  '[' constant_expression ']' '=' constant_expression
;

compound_statement:
  '{' statement_list '}'
;

struct_body:
  '{' struct_contents '}'
;

union_body:
  '{' struct_contents '}'
;

enum_body:
  '{' enum_element_list '}'
;

statement_list:
  /*empty */
| statement_list statement
| statement_list external_declaration  // allows local scopes to define, types, functions, etc
;

struct_contents:
  struct_element_list
| struct_element_list ELIPSIS
;

struct_element_list:
  /* empty */
| struct_element_list struct_element
;

struct_element:
  declaration_type IDENTIFIER variable_properties initializer ';'
;

enum_element_list:
  /* empty */
| enum_element
;

enum_element:
  IDENTIFIER ';'
| IDENTIFIER '=' constant_expression ';'
;

external_declaration:
  declaration_type IDENTIFIER variable_properties initializer ';'
| declaration_type IDENTIFIER '(' proto_params ')' function_properties ';'
| declaration_type IDENTIFIER '(' proto_params ')' function_properties compound_statement
| STRUCT IDENTIFIER struct_properties ';'
| STRUCT IDENTIFIER struct_properties struct_body
| UNION IDENTIFIER union_properties ';'
| UNION IDENTIFIER union_properties union_body
| ENUM type_expression IDENTIFIER enum_properties ';'
| ENUM type_expression IDENTIFIER enum_properties enum_body
;

statement:
  compound_statement
| expression ';'
| label_statement
| for_statement
| while_statement
| do_statement
| if_statement
| switch_statement
| break_statement
| continue_statement
| goto_statement
| return_statement
| ';'
;

label_statement:
  IDENTIFIER ':' statement
;

for_statement:
  FOR '(' expression ';' expression ';' expression ')' statement
;

while_statement:
  WHILE '(' expression ')' statement
;

do_statement:
  DO statement WHILE '(' expression ')' ';'
;

if_statement:
  IF '(' expression ')' statement
| IF '(' expression ')' statement ELSE statement //ambiguous, use greedy parse
;

switch_statement:
  SWITCH '(' expression ')' '{' cases '}'
;

break_statement:
  BREAK ';'
  // FIXME think about other break rules
;

continue_statement:
  CONTINUE ';'
  // FIXME think about nested loop continue rules
;

goto_statement:
  GOTO IDENTIFIER ';'
// FIXME think about goto rules
;

return_statement:
  RETURN ';'
| RETURN expression ';'
;

cases:
  /* empty */
| case_list case_last_element
;

case_list:
  /* empty */
| case_list case_element
;

case_element:
  case_start_statement statement_list case_end_statement
;

case_last_element:
  case_start_statement statement_list
;

case_start_statement:
  DEFAULT ':'
| CASE IDENTIFIER ':'
;

case_end_statement:
  break_statement
| return_statement
| continue_statement
| goto_statement
;

expressions:
  /* empty */
| expression_list
| expression_list ',' // extra trailing comma allowed
;

expression_list:
  conditional_expression
| expression_list ',' conditional_expression
;

constant_expression:
  conditional_expression
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
| and_expression AND_ADDR_OP equality_expression
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
| multiplicative_expression DIV_OP unary_expression
| multiplicative_expression MULT_PTR_OP unary_expression
  //no cast expression, casts look like function calls now
;

unary_expression:
  postfix_expression
| SIZEOF unary_expression
| NOT_OP unary_expression
| MULT_PTR_OP unary_expression
| AND_ADDR_OP unary_expression
| ADD_OP unary_expression
| INC_OP unary_expression
;

postfix_expression:
  primary_expression
| postfix_expression '[' expressions ']'
| postfix_expression '(' expressions ')'
| postfix_expression '{' initializers '}'
| postfix_expression '.' IDENTIFIER // ambiguous with namespace_expression
| postfix_expression INC_OP
// no arrow operator postfix_expression -> IDENTIFIER
;

primary_expression:
  type_expression
| NUMBER
| CHARCONST
| string_const
| '(' expression ')'
;

string_const:
  STRINGCONST
| string_const STRINGCONST
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
| READONLY
| VOLATILE
| ARRAY '(' expressions ')' // FIXME does this need to be constant ?
| ARRAY '(' '*' ')'
| BITFIELD '(' constant_expression ')'
//FIXME maybe noaddress to prevent making pointers to it
;

storage_specifier:
  CONSTANT
| TYPEDEF
| ALIAS
| REGISTER
| STATIC    // not allowed at file scope
| AUTO
;

parameter_specifier:
  RESTRICT
| REGISTER
;

variable_properties:
 /* empty */
| LINKAGE '(' string_const ')' // FIXME maybe these are any constant expression
| LINKNAME '(' string_const ')'
// extern("image.jpg", "binary")
// extern("helpfile.txt", "utf-8") // adds zero terminator to string
// FIXME more...
;

function_properties:
 /* empty */
| LINKAGE '(' string_const ')' // FIXME maybe these are any constant expression
| LINKNAME '(' string_const ')'
| INLINE
// FIXME more...
;

struct_properties:
  /* empty */
| ALLIGNMENT
// FIXME more ...
;

union_properties:
  /* empty */
| ALLIGNMENT
// FIXME more ...
;

enum_properties:
  /* empty */
| ALLIGNMENT
// size(constant_expression)
// FIXME more ...
;

%%

int yylex(void)
{
	yylval = 1234;
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

double MakeNode(int count, double *params)
{
	return 0;
}

