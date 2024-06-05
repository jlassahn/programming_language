
%{
#include <stdio.h>
#include "compiler/parser.h"

// Bison C interface adjustments
int yylex(void);
void yyerror(const char *s);

// yylen is an internal symbol that counts the number of inputs to the rule.
// yyvsp is the internal stack holding semantic values.
#define MkNode(kind) ((yyval) = MakeNode(kind, yylen, &yyvsp[1-yylen]))
#define MkEmpty ((yyval) = MakeNode(&SYM_EMPTY, 0, NULL))
#define MkMove ((yyval) = yyvsp[1-yylen])

#define YYSTYPE ParserNode *

%}

/* punctuators
  ; { } , . [ ] ( ) :
   conditional op
  ? :
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
%token ASSIGN_OP      /* = */
%token ASSIGN_MULT_OP /* *= */
%token ASSIGN_DIV_OP  /* /= */
%token ASSIGN_MOD_OP  /* %= */
%token ASSIGN_ADD_OP  /* += */
%token ASSIGN_SUB_OP  /* -= */
%token ASSIGN_SHR_OP  /* >>= */
%token ASSIGN_SHL_OP  /* <<= */
%token ASSIGN_AND_OP  /* &= */
%token ASSIGN_OR_OP   /* |= */
%token ASSIGN_XOR_OP  /* ^= */
%token LOG_OR_OP      /* || */
%token LOG_AND_OP     /* && */
%token OR_OP          /* | */
%token AND_ADDR_OP    /* & */
%token XOR_OP         /* ^ */
%token EQUAL_OP       /* == */
%token NEQUAL_OP      /* != */
%token LESS_OP        /* < */
%token GREATER_OP     /* > */
%token LESSEQ_OP      /* <= */
%token GREATEREQ_OP   /* >= */
%token SHL_OP         /* << */
%token SHR_OP         /* >> */
%token ADD_OP         /* + */
%token SUB_OP         /* - */
%token DIV_OP         /* / */
%token MOD_OP         /* % */
%token MULT_PTR_OP    /* * */
%token NOT_OP         /* ! */
%token BITNOT_OP      /* ~ */
%token INC_OP         /* ++ */
%token DEC_OP         /* -- */

%token ELIPSIS

%%

file:
  /* empty */       { MkEmpty; }
| file file_element { MkNode(&SYM_FIXME); }
// FIXME maybe a file header?
;

file_element:
  import_statement     { MkMove; }
| using_statement      { MkMove; }
| external_declaration { MkMove; }
;

import_statement:
  IMPORT namespace_expression ';'         { MkNode(&SYM_FIXME); }
| IMPORT PRIVATE namespace_expression ';' { MkNode(&SYM_FIXME); }
;

using_statement:
  USING namespace_expression ';'                         { MkNode(&SYM_FIXME); }
| USING namespace_expression AS namespace_expression ';' { MkNode(&SYM_FIXME); }
;

// extra trailing comma allowed
proto_params:
  /* empty */                  { MkEmpty; }
| proto_param_list             { MkMove; }
| proto_param_list ','         { MkNode(&SYM_FIXME); }
| proto_param_list ',' ELIPSIS { MkNode(&SYM_FIXME); }
;

proto_param_list:
  proto_param                      { MkMove; }
| proto_param_list ',' proto_param { MkNode(&SYM_FIXME); }
;

proto_param:
  param_type                        { MkMove; }
| param_type IDENTIFIER initializer { MkNode(&SYM_FIXME); }
;

param_type:
  type_expression                     { MkMove; }
| parameter_specifier type_expression { MkNode(&SYM_FIXME); }
;

declaration_type:
  type_expression                   { MkMove; }
| storage_specifier type_expression { MkNode(&SYM_FIXME); }
;

initializer:
  /* empty */                    { MkEmpty; }
| ASSIGN_OP expression           { MkNode(&SYM_FIXME); }
| ASSIGN_OP '{' initializers '}' { MkNode(&SYM_FIXME); }
;

// extra trailing comma allowed
initializers:
  /* empty */           { MkEmpty; }
| initializer_list      { MkMove; }
| initializer_list ','  { MkNode(&SYM_FIXME); }
;

initializer_list:
  initializer_element                      { MkMove; }
| initializer_list ',' initializer_element { MkNode(&SYM_FIXME); }
;

initializer_element:
  '{' initializers '}' { MkNode(&SYM_FIXME); }
| struct_initializer   { MkNode(&SYM_FIXME); }
| array_initializer    { MkNode(&SYM_FIXME); }
| constant_expression  { MkNode(&SYM_FIXME); }
;

struct_initializer:
  '.' IDENTIFIER ASSIGN_OP constant_expression { MkNode(&SYM_FIXME); }
;

array_initializer:
  '[' constant_expression ']' ASSIGN_OP constant_expression { MkNode(&SYM_FIXME); }
;

compound_statement:
  '{' statement_list '}' { MkNode(&SYM_FIXME); }
;

struct_body:
  '{' struct_contents '}' { MkNode(&SYM_FIXME); }
;

union_body:
  '{' struct_contents '}' { MkNode(&SYM_FIXME); }
;

enum_body:
  '{' enum_element_list '}' { MkNode(&SYM_FIXME); }
;

// allows local scopes to define, types, functions, etc
statement_list:
  /*empty */                           { MkEmpty; }
| statement_list statement             { MkNode(&SYM_FIXME); }
| statement_list external_declaration  { MkNode(&SYM_FIXME); }
;

struct_contents:
  struct_element_list         { MkMove; }
| struct_element_list ELIPSIS { MkNode(&SYM_FIXME); }
;

struct_element_list:
  /* empty */                        { MkEmpty; }
| struct_element_list struct_element { MkNode(&SYM_FIXME); }
;

struct_element:
  declaration_type IDENTIFIER variable_properties initializer ';' { MkNode(&SYM_FIXME); }
;

enum_element_list:
  /* empty */                    { MkEmpty; }
| enum_element_list enum_element { MkNode(&SYM_FIXME); }
;

enum_element:
  IDENTIFIER ';'                               { MkNode(&SYM_FIXME); }
| IDENTIFIER ASSIGN_OP constant_expression ';' { MkNode(&SYM_FIXME); }
;

external_declaration:
  declaration_type IDENTIFIER variable_properties initializer ';'                         { MkNode(&SYM_FIXME); }
| declaration_type IDENTIFIER '(' proto_params ')' function_properties ';'                { MkNode(&SYM_FIXME); }
| declaration_type IDENTIFIER '(' proto_params ')' function_properties compound_statement { MkNode(&SYM_FIXME); }
| STRUCT IDENTIFIER struct_properties ';'                                                 { MkNode(&SYM_FIXME); }
| STRUCT IDENTIFIER struct_properties struct_body                                         { MkNode(&SYM_FIXME); }
| UNION IDENTIFIER union_properties ';'                                                   { MkNode(&SYM_FIXME); }
| UNION IDENTIFIER union_properties union_body                                            { MkNode(&SYM_FIXME); }
| ENUM type_expression IDENTIFIER enum_properties ';'                                     { MkNode(&SYM_FIXME); }
| ENUM type_expression IDENTIFIER enum_properties enum_body                               { MkNode(&SYM_FIXME); }
;

statement:
  compound_statement { MkMove; }
| expression ';'     { MkNode(&SYM_FIXME); }
| label_statement    { MkMove; }
| for_statement      { MkMove; }
| while_statement    { MkMove; }
| do_statement       { MkMove; }
| if_statement       { MkMove; }
| switch_statement   { MkMove; }
| break_statement    { MkMove; }
| continue_statement { MkMove; }
| goto_statement     { MkMove; }
| return_statement   { MkMove; }
| ';'                { MkNode(&SYM_FIXME); }
;

label_statement:
  IDENTIFIER ':' statement { MkNode(&SYM_FIXME); }
;

for_statement:
  FOR '(' for_initializer ';' expression ';' expression ')' statement { MkNode(&SYM_FIXME); }
;

for_initializer:
  expression                                                  { MkMove; }
| declaration_type IDENTIFIER variable_properties initializer { MkNode(&SYM_FIXME); }
;

while_statement:
  WHILE '(' expression ')' statement { MkNode(&SYM_FIXME); }
;

do_statement:
  DO statement WHILE '(' expression ')' ';' { MkNode(&SYM_FIXME); }
;

//ambiguous, use greedy parse
if_statement:
  IF '(' expression ')' statement                 { MkNode(&SYM_FIXME); }
| IF '(' expression ')' statement ELSE statement  { MkNode(&SYM_FIXME); }
;

switch_statement:
  SWITCH '(' expression ')' '{' cases '}' { MkNode(&SYM_FIXME); }
;

break_statement:
  BREAK ';' { MkNode(&SYM_FIXME); }
  // FIXME think about other break rules
;

continue_statement:
  CONTINUE ';' { MkNode(&SYM_FIXME); }
  // FIXME think about nested loop continue rules
;

goto_statement:
  GOTO IDENTIFIER ';' { MkNode(&SYM_FIXME); }
// FIXME think about goto rules
;

return_statement:
  RETURN ';' { MkNode(&SYM_FIXME); }
| RETURN expression ';' { MkNode(&SYM_FIXME); }
;

cases:
  /* empty */                 { MkEmpty; }
| case_list case_last_element { MkNode(&SYM_FIXME); }
;

case_list:
  /* empty */            { MkEmpty; }
| case_list case_element { MkNode(&SYM_FIXME); }
;

case_element:
  case_start_statement statement_list case_end_statement { MkNode(&SYM_FIXME); }
;

case_last_element:
  case_start_statement statement_list { MkNode(&SYM_FIXME); }
;

case_start_statement:
  DEFAULT ':'         { MkNode(&SYM_FIXME); }
| CASE IDENTIFIER ':' { MkNode(&SYM_FIXME); }
;

case_end_statement:
  break_statement    { MkMove; }
| return_statement   { MkMove; }
| continue_statement { MkMove; }
| goto_statement     { MkMove; }
;

// extra trailing comma allowed
expressions:
  /* empty */          { MkEmpty; }
| expression_list      { MkMove; }
| expression_list ','  { MkNode(&SYM_FIXME); }
;

expression_list:
  conditional_expression                     { MkMove; }
| expression_list ',' conditional_expression { MkNode(&SYM_FIXME); }
;

constant_expression:
  conditional_expression { MkNode(&SYM_FIXME); }
;

expression:
  assignment_expression { MkMove; }
// no comma operator
;

assignment_expression:
  conditional_expression                               { MkMove; }
| unary_expression assignment_op assignment_expression { MkNode(&SYM_FIXME); }
;

assignment_op:
  ASSIGN_OP      { MkMove; }
| ASSIGN_MULT_OP { MkMove; }
| ASSIGN_DIV_OP  { MkMove; }
| ASSIGN_MOD_OP  { MkMove; }
| ASSIGN_ADD_OP  { MkMove; }
| ASSIGN_SUB_OP  { MkMove; }
| ASSIGN_SHR_OP  { MkMove; }
| ASSIGN_SHL_OP  { MkMove; }
| ASSIGN_AND_OP  { MkMove; }
| ASSIGN_OR_OP   { MkMove; }
| ASSIGN_XOR_OP  { MkMove; }
;

conditional_expression:
  logical_or_expression                                           { MkMove; }
| logical_or_expression '?' expression ':' conditional_expression { MkNode(&SYM_FIXME); }
;

logical_or_expression:
  logical_and_expression                                 { MkMove; }
| logical_or_expression LOG_OR_OP logical_and_expression { MkNode(&SYM_FIXME); }
;

logical_and_expression:
  inclusive_or_expression                                   { MkMove; }
| logical_and_expression LOG_AND_OP inclusive_or_expression { MkNode(&SYM_FIXME); }
;

inclusive_or_expression:
  exclusive_or_expression                               { MkMove; }
| inclusive_or_expression OR_OP exclusive_or_expression { MkNode(&SYM_FIXME); }
;

exclusive_or_expression:
  and_expression                                { MkMove; }
| exclusive_or_expression XOR_OP and_expression { MkNode(&SYM_FIXME); }
;

and_expression:
  equality_expression                            { MkMove; }
| and_expression AND_ADDR_OP equality_expression { MkNode(&SYM_FIXME); }
;

equality_expression:
  relational_expression                               { MkMove; }
| equality_expression EQUAL_OP relational_expression  { MkNode(&SYM_FIXME); }
| equality_expression NEQUAL_OP relational_expression { MkNode(&SYM_FIXME); }
;

relational_expression:
  shift_expression                                    { MkMove; }
| relational_expression LESS_OP shift_expression      { MkNode(&SYM_FIXME); }
| relational_expression GREATER_OP shift_expression   { MkNode(&SYM_FIXME); }
| relational_expression LESSEQ_OP shift_expression    { MkNode(&SYM_FIXME); }
| relational_expression GREATEREQ_OP shift_expression { MkNode(&SYM_FIXME); }
;

shift_expression:
  additive_expression                         { MkMove; }
| shift_expression SHL_OP additive_expression { MkNode(&SYM_FIXME); }
| shift_expression SHR_OP additive_expression { MkNode(&SYM_FIXME); }
;

additive_expression:
  multiplicative_expression                            { MkMove; }
| additive_expression ADD_OP multiplicative_expression { MkNode(&SYM_FIXME); }
| additive_expression SUB_OP multiplicative_expression { MkNode(&SYM_FIXME); }
;

multiplicative_expression:
  unary_expression                                  { MkMove; }
| multiplicative_expression DIV_OP unary_expression { MkNode(&SYM_FIXME); }
| multiplicative_expression MOD_OP unary_expression { MkNode(&SYM_FIXME); }
| multiplicative_expression MULT_PTR_OP unary_expression { MkNode(&SYM_FIXME); }
  //no cast expression, casts look like function calls now
;

unary_expression:
  postfix_expression           { MkMove; }
| SIZEOF unary_expression      { MkNode(&SYM_FIXME); }
| NOT_OP unary_expression      { MkNode(&SYM_FIXME); }
| BITNOT_OP unary_expression   { MkNode(&SYM_FIXME); }
| MULT_PTR_OP unary_expression { MkNode(&SYM_FIXME); }
| AND_ADDR_OP unary_expression { MkNode(&SYM_FIXME); }
| ADD_OP unary_expression      { MkNode(&SYM_FIXME); }
| SUB_OP unary_expression      { MkNode(&SYM_FIXME); }
| INC_OP unary_expression      { MkNode(&SYM_FIXME); }
| DEC_OP unary_expression      { MkNode(&SYM_FIXME); }
;

// dot operator ambiguous with namespace_expression
postfix_expression:
  primary_expression                      { MkMove; }
| postfix_expression '[' expressions ']'  { MkNode(&SYM_FIXME); }
| postfix_expression '(' expressions ')'  { MkNode(&SYM_FIXME); }
| postfix_expression '{' initializers '}' { MkNode(&SYM_FIXME); }
| postfix_expression '.' IDENTIFIER       { MkNode(&SYM_FIXME); }
| postfix_expression INC_OP               { MkNode(&SYM_FIXME); }
| postfix_expression DEC_OP               { MkNode(&SYM_FIXME); }
// no arrow operator postfix_expression -> IDENTIFIER
;

primary_expression:
  type_expression    { MkMove; }
| NUMBER             { MkMove; }
| CHARCONST          { MkMove; }
| string_const       { MkMove; }
| '(' expression ')' { MkNode(&SYM_FIXME); }
;

string_const:
  STRINGCONST              { MkMove; }
| string_const STRINGCONST { MkNode(&SYM_FIXME); }
;

type_expression:
  namespace_expression          { MkMove; }
| type_modifier type_expression { MkNode(&SYM_FIXME); }
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

// dot operator ambiguous with postfix_expression
namespace_expression:
  IDENTIFIER                           { MkMove; }
| namespace_expression '.' IDENTIFIER  { MkNode(&SYM_FIXME); }
;

// FIXME do array sizes need to be constant ?
type_modifier:
  POINTER                              { MkMove; }
| READONLY                             { MkMove; }
| VOLATILE                             { MkMove; }
| ARRAY '(' expressions ')'            { MkNode(&SYM_FIXME); }
| ARRAY '(' MULT_PTR_OP ')'            { MkNode(&SYM_FIXME); }
| BITFIELD '(' constant_expression ')' { MkNode(&SYM_FIXME); }
//FIXME maybe noaddress to prevent making pointers to it
;

storage_specifier:
  CONSTANT { MkMove; }
| TYPEDEF  { MkMove; }
| ALIAS    { MkMove; }
| REGISTER { MkMove; }
| STATIC   { MkMove; }  // not allowed at file scope
| AUTO     { MkMove; }
;

parameter_specifier:
  RESTRICT { MkMove; }
| REGISTER { MkMove; }
;

// FIXME maybe parameters are any constant expression
variable_properties:
 /* empty */                    { MkEmpty; }
| LINKAGE '(' string_const ')'  { MkNode(&SYM_FIXME); }
| LINKNAME '(' string_const ')' { MkNode(&SYM_FIXME); }
// extern("image.jpg", "binary")
// extern("helpfile.txt", "utf-8") // adds zero terminator to string
// FIXME more...
;

// FIXME maybe parameters are any constant expression
function_properties:
 /* empty */                    { MkEmpty; }
| LINKAGE '(' string_const ')'  { MkNode(&SYM_FIXME); }
| LINKNAME '(' string_const ')' { MkNode(&SYM_FIXME); }
| INLINE                        { MkMove; }
// FIXME more...
;

struct_properties:
  /* empty */  { MkEmpty; }
| ALLIGNMENT   { MkMove; }
// FIXME more ...
;

union_properties:
  /* empty */  { MkEmpty; }
| ALLIGNMENT   { MkMove; }
// FIXME more ...
;

enum_properties:
  /* empty */  { MkEmpty; }
| ALLIGNMENT   { MkMove; }
// size(constant_expression)
// FIXME more ...
;

%%

