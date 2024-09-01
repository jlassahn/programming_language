
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
%token FUNCTION
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

%token ELLIPSIS

%%

file:
  /* empty */       { MkEmpty; }
| file file_element { MkNode(&SYM_LIST); }
// FIXME maybe a file header?
;

file_element:
  import_statement     { MkMove; }
| using_statement      { MkMove; }
| external_declaration { MkMove; }
| ';'                  { MkNode(&SYM_EMPTY_STATEMENT); }
;

import_statement:
  IMPORT namespace_expression ';'         { MkNode(&SYM_IMPORT); }
| IMPORT PRIVATE namespace_expression ';' { MkNode(&SYM_IMPORT_PRIVATE); }
;

using_statement:
  USING namespace_expression ';'                         { MkNode(&SYM_USING); }
| USING namespace_expression AS namespace_expression ';' { MkNode(&SYM_USING_AS); }
;

// extra trailing comma allowed
proto_params:
  /* empty */                   { MkEmpty; }
| proto_param_list              { MkMove; }
| proto_param_list ','          { MkNode(&SYM_TRAILING_COMMA); }
| proto_param_list ',' ELLIPSIS { MkNode(&SYM_ELLIPSIS); }
;

proto_param_list:
  proto_param                      { MkMove; }
| proto_param_list ',' proto_param { MkNode(&SYM_LIST); }
;

proto_param:
  param_type                        { MkMove; }
| param_type IDENTIFIER initializer { MkNode(&SYM_PARAMETER); }
;

param_type:
  type_expression                     { MkMove; }
| parameter_specifier type_expression { MkNode(&SYM_PARAM_TYPE); }
;

declaration_type:
  type_expression                   { MkMove; }
| storage_specifier type_expression { MkNode(&SYM_DECL_TYPE); }
;

initializer:
  /* empty */                    { MkEmpty; }
| ASSIGN_OP expression           { MkNode(&SYM_INITIALIZE); }
| ASSIGN_OP '{' initializers '}' { MkNode(&SYM_INITIALIZE); }
;

// extra trailing comma allowed
initializers:
  /* empty */           { MkEmpty; }
| initializer_list      { MkMove; }
| initializer_list ','  { MkNode(&SYM_TRAILING_COMMA); }
;

initializer_list:
  initializer_element                      { MkMove; }
| initializer_list ',' initializer_element { MkNode(&SYM_LIST); }
;

initializer_element:
  '{' initializers '}' { MkNode(&SYM_INITIALIZE); }
| struct_initializer   { MkMove; }
| array_initializer    { MkMove; }
| constant_expression  { MkMove; }
;

struct_initializer:
  '.' IDENTIFIER ASSIGN_OP constant_expression { MkNode(&SYM_INIT_STRUCT); }
;

array_initializer:
  '[' constant_expression ']' ASSIGN_OP constant_expression { MkNode(&SYM_INIT_ARRAY); }
;

compound_statement:
  '{' statement_list '}' { MkNode(&SYM_STATEMENT_LIST); }
;

struct_body:
  '{' struct_contents '}' { MkNode(&SYM_STRUCT_LIST); }
;

union_body:
  '{' struct_contents '}' { MkNode(&SYM_STRUCT_LIST); }
;

enum_body:
  '{' enum_element_list '}' { MkNode(&SYM_ENUM_LIST); }
;

// allows local scopes to define, types, functions, etc
statement_list:
  /*empty */                           { MkEmpty; }
| statement_list statement             { MkNode(&SYM_LIST); }
| statement_list external_declaration  { MkNode(&SYM_LIST); }
;

struct_contents:
  struct_element_list          { MkMove; }
| struct_element_list ELLIPSIS { MkNode(&SYM_ELLIPSIS); }
;

struct_element_list:
  /* empty */                        { MkEmpty; }
| struct_element_list struct_element { MkNode(&SYM_LIST); }
;

struct_element:
  declaration_type IDENTIFIER variable_properties initializer ';' { MkNode(&SYM_DECLARATION); }
;

enum_element_list:
  /* empty */                    { MkEmpty; }
| enum_element_list enum_element { MkNode(&SYM_LIST); }
;

enum_element:
  IDENTIFIER initializer ';'  { MkNode(&SYM_ENUM_ELEMENT); }
;

external_declaration:
  declaration_type IDENTIFIER variable_properties initializer ';'                         { MkNode(&SYM_DECLARATION); }
| declaration_type IDENTIFIER '(' proto_params ')' function_properties ';'                { MkNode(&SYM_PROTOTYPE); }
| declaration_type IDENTIFIER '(' proto_params ')' function_properties compound_statement { MkNode(&SYM_FUNC); }
| STRUCT IDENTIFIER struct_properties ';'                                                 { MkNode(&SYM_STRUCT_DEC); }
| STRUCT IDENTIFIER struct_properties struct_body                                         { MkNode(&SYM_STRUCT_DEF); }
| UNION IDENTIFIER union_properties ';'                                                   { MkNode(&SYM_UNION_DEC); }
| UNION IDENTIFIER union_properties union_body                                            { MkNode(&SYM_UNION_DEF); }
| ENUM type_expression IDENTIFIER enum_properties ';'                                     { MkNode(&SYM_ENUM_DEC); }
| ENUM type_expression IDENTIFIER enum_properties enum_body                               { MkNode(&SYM_ENUM_DEF); }
;

statement:
  compound_statement { MkMove; }
| expression ';'     { MkNode(&SYM_EXPRESSION_STATEMENT); }
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
| ';'                { MkNode(&SYM_EMPTY_STATEMENT); }
;

label_statement:
  IDENTIFIER ':' statement { MkNode(&SYM_LABEL_STATEMENT); }
;

for_statement:
  FOR '(' for_initializer ';' expression ';' expression ')' statement { MkNode(&SYM_FOR_STATEMENT); }
;

for_initializer:
  expression                                                  { MkMove; }
| declaration_type IDENTIFIER variable_properties initializer { MkNode(&SYM_DECLARATION); }
;

while_statement:
  WHILE '(' expression ')' statement { MkNode(&SYM_WHILE_STATEMENT); }
;

do_statement:
  DO statement WHILE '(' expression ')' ';' { MkNode(&SYM_DO_STATEMENT); }
;

//ambiguous, use greedy parse
if_statement:
  IF '(' expression ')' statement                 { MkNode(&SYM_IF_STATEMENT); }
| IF '(' expression ')' statement ELSE statement  { MkNode(&SYM_IF_ELSE); }
;

switch_statement:
  SWITCH '(' expression ')' '{' cases '}' { MkNode(&SYM_SWITCH_STATEMENT); }
;

break_statement:
  BREAK ';' { MkNode(&SYM_BREAK_STATEMENT); }
  // FIXME think about other break rules
;

continue_statement:
  CONTINUE ';' { MkNode(&SYM_CONTINUE_STATEMENT); }
  // FIXME think about nested loop continue rules
;

goto_statement:
  GOTO IDENTIFIER ';' { MkNode(&SYM_GOTO_STATEMENT); }
// FIXME think about goto rules
;

return_statement:
  RETURN ';' { MkNode(&SYM_RETURN_VOID); }
| RETURN expression ';' { MkNode(&SYM_RETURN_STATEMENT); }
;

cases:
  /* empty */                 { MkEmpty; }
| case_list case_last_element { MkNode(&SYM_LIST); }
;

case_list:
  /* empty */            { MkEmpty; }
| case_list case_element { MkNode(&SYM_LIST); }
;

case_element:
  case_start statement_list case_end_statement { MkNode(&SYM_CASE_ELEMENT); }
;

case_last_element:
  case_start statement_list { MkNode(&SYM_CASE_END_ELEMENT); }
;

case_start:
  case_start_statement            { MkMove; }
| case_start case_start_statement { MkNode(&SYM_LIST); }
;

case_start_statement:
  DEFAULT ':'                  { MkNode(&SYM_DEFAULT_LABEL); }
| CASE constant_expression ':' { MkNode(&SYM_CASE_LABEL); }
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
| expression_list ','  { MkNode(&SYM_TRAILING_COMMA); }
;

expression_list:
  conditional_expression                     { MkMove; }
| expression_list ',' conditional_expression { MkNode(&SYM_LIST); }
;

constant_expression:
  conditional_expression { MkNode(&SYM_CONSTANT); }
;

expression:
  assignment_expression { MkMove; }
// no comma operator
;

assignment_expression:
  conditional_expression                                { MkMove; }
| unary_expression ASSIGN_OP assignment_expression      { MkNode(&SYM_ASSIGN_OP); }
| unary_expression ASSIGN_MULT_OP assignment_expression { MkNode(&SYM_ASSIGN_MULT_OP); }
| unary_expression ASSIGN_DIV_OP assignment_expression  { MkNode(&SYM_ASSIGN_DIV_OP); }
| unary_expression ASSIGN_MOD_OP assignment_expression  { MkNode(&SYM_ASSIGN_MOD_OP); }
| unary_expression ASSIGN_ADD_OP assignment_expression  { MkNode(&SYM_ASSIGN_ADD_OP); }
| unary_expression ASSIGN_SUB_OP assignment_expression  { MkNode(&SYM_ASSIGN_SUB_OP); }
| unary_expression ASSIGN_SHR_OP assignment_expression  { MkNode(&SYM_ASSIGN_SHR_OP); }
| unary_expression ASSIGN_SHL_OP assignment_expression  { MkNode(&SYM_ASSIGN_SHL_OP); }
| unary_expression ASSIGN_AND_OP assignment_expression  { MkNode(&SYM_ASSIGN_AND_OP); }
| unary_expression ASSIGN_OR_OP assignment_expression   { MkNode(&SYM_ASSIGN_OR_OP); }
| unary_expression ASSIGN_XOR_OP assignment_expression  { MkNode(&SYM_ASSIGN_XOR_OP); }
;

conditional_expression:
  logical_or_expression                                           { MkMove; }
| logical_or_expression '?' expression ':' conditional_expression { MkNode(&SYM_CONDITIONAL); }
;

logical_or_expression:
  logical_and_expression                                 { MkMove; }
| logical_or_expression LOG_OR_OP logical_and_expression { MkNode(&SYM_LOG_OR_OP); }
;

logical_and_expression:
  inclusive_or_expression                                   { MkMove; }
| logical_and_expression LOG_AND_OP inclusive_or_expression { MkNode(&SYM_LOG_AND_OP); }
;

inclusive_or_expression:
  exclusive_or_expression                               { MkMove; }
| inclusive_or_expression OR_OP exclusive_or_expression { MkNode(&SYM_OR_OP); }
;

exclusive_or_expression:
  and_expression                                { MkMove; }
| exclusive_or_expression XOR_OP and_expression { MkNode(&SYM_XOR_OP); }
;

and_expression:
  equality_expression                            { MkMove; }
| and_expression AND_ADDR_OP equality_expression { MkNode(&SYM_AND_OP); }
;

equality_expression:
  relational_expression                               { MkMove; }
| equality_expression EQUAL_OP relational_expression  { MkNode(&SYM_EQUAL_OP); }
| equality_expression NEQUAL_OP relational_expression { MkNode(&SYM_NEQUAL_OP); }
;

relational_expression:
  shift_expression                                    { MkMove; }
| relational_expression LESS_OP shift_expression      { MkNode(&SYM_LESS_OP); }
| relational_expression GREATER_OP shift_expression   { MkNode(&SYM_GREATER_OP); }
| relational_expression LESSEQ_OP shift_expression    { MkNode(&SYM_LESSEQ_OP); }
| relational_expression GREATEREQ_OP shift_expression { MkNode(&SYM_GREATEREQ_OP); }
;

shift_expression:
  additive_expression                         { MkMove; }
| shift_expression SHL_OP additive_expression { MkNode(&SYM_SHL_OP); }
| shift_expression SHR_OP additive_expression { MkNode(&SYM_SHR_OP); }
;

additive_expression:
  multiplicative_expression                            { MkMove; }
| additive_expression ADD_OP multiplicative_expression { MkNode(&SYM_ADD_OP); }
| additive_expression SUB_OP multiplicative_expression { MkNode(&SYM_SUB_OP); }
;

multiplicative_expression:
  unary_expression                                  { MkMove; }
| multiplicative_expression DIV_OP unary_expression { MkNode(&SYM_DIV_OP); }
| multiplicative_expression MOD_OP unary_expression { MkNode(&SYM_MOD_OP); }
| multiplicative_expression MULT_PTR_OP unary_expression { MkNode(&SYM_MULT_OP); }
  //no cast expression, casts look like function calls now
;

unary_expression:
  postfix_expression           { MkMove; }
| SIZEOF unary_expression      { MkNode(&SYM_SIZEOF_OP); }
| NOT_OP unary_expression      { MkNode(&SYM_NOT_OP); }
| BITNOT_OP unary_expression   { MkNode(&SYM_BITNOT_OP); }
| MULT_PTR_OP unary_expression { MkNode(&SYM_PTR_OP); }
| AND_ADDR_OP unary_expression { MkNode(&SYM_ADDR_OP); }
| ADD_OP unary_expression      { MkNode(&SYM_POS_OP); }
| SUB_OP unary_expression      { MkNode(&SYM_NEG_OP); }
| INC_OP unary_expression      { MkNode(&SYM_PREINC_OP); }
| DEC_OP unary_expression      { MkNode(&SYM_PREDEC_OP); }
;

// dot operator ambiguous with namespace_expression
postfix_expression:
  primary_expression                      { MkMove; }
| postfix_expression '[' expressions ']'  { MkNode(&SYM_ARRAY_OP); }
| postfix_expression '(' expressions ')'  { MkNode(&SYM_CALL_OP); }
| postfix_expression '{' initializers '}' { MkNode(&SYM_INIT_OP); }
| postfix_expression '.' IDENTIFIER       { MkNode(&SYM_DOT_OP); }
| postfix_expression INC_OP               { MkNode(&SYM_POSTINC_OP); }
| postfix_expression DEC_OP               { MkNode(&SYM_POSTDEC_OP); }
// no arrow operator postfix_expression -> IDENTIFIER
;

primary_expression:
  type_expression    { MkMove; }
| NUMBER             { MkMove; }
| CHARCONST          { MkMove; }
| string_const       { MkMove; }
| '(' expression ')' { MkNode(&SYM_PAREN_EXPRESSION); }
;

string_const:
  STRINGCONST              { MkMove; }
| string_const STRINGCONST { MkNode(&SYM_STRING); }
;

type_expression:
  namespace_expression          { MkMove; }
| type_modifier type_expression { MkNode(&SYM_TYPE_EXPRESSION); }
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
| namespace_expression '.' IDENTIFIER  { MkNode(&SYM_DOT_OP); }
;

// FIXME do array sizes need to be constant ?
type_modifier:
  POINTER                              { MkMove; }
| READONLY                             { MkMove; }
| VOLATILE                             { MkMove; }
| ARRAY '(' expressions ')'            { MkNode(&SYM_TYPE_ARRAY); }
| ARRAY '(' MULT_PTR_OP ')'            { MkNode(&SYM_TYPE_ARRAY_MATCH); }
| BITFIELD '(' constant_expression ')' { MkNode(&SYM_TYPE_BITFIELD); }
| FUNCTION '(' proto_params ')'        { MkNode(&SYM_TYPE_FUNCTION); }
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
| LINKAGE '(' string_const ')'  { MkNode(&SYM_TYPE_LINKAGE); }
| LINKNAME '(' string_const ')' { MkNode(&SYM_TYPE_LINKNAME); }
// extern("image.jpg", "binary")
// extern("helpfile.txt", "utf-8") // adds zero terminator to string
// extern("DEBUG", "option") // allows compile-time options
// noinit
// FIXME more...
;

// FIXME maybe parameters are any constant expression
function_properties:
 /* empty */                    { MkEmpty; }
| LINKAGE '(' string_const ')'  { MkNode(&SYM_TYPE_LINKAGE); }
| LINKNAME '(' string_const ')' { MkNode(&SYM_TYPE_LINKNAME); }
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

