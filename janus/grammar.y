
%{
%}

%token NUMBER_TOKEN
%token STRING_TOKEN
%token CHARACTER_TOKEN
%token SYMBOL_TOKEN

%token OR_OP
%token AND_OP
%token COMPARE_OP
%token ADD_OP
%token MULT_OP
%token EXP_OP
%token PREFIX_OP
%token SUFFIX_OP
%token ASSIGNMENT_OP
%token ANY_OP

%token JANUS
%token IMPORT
%token FUNCTION
%token DEF
%token CONST
%token STRUCT_OR_MSTRUCT
%token EXTENDS
%token IMPLEMENTS
/* FIXME do we want???  %token SIZE */
%token ALIAS
%token INTERFACE
%token METHOD
%token OPERATOR
%token IF
%token ELSE
%token WHILE
%token FOR
%token WITH
%token RETURN
%token BREAK
%token CONTINUE
%token LABEL
%token GOTO


%%

file:
	header
	| file  file_declaration
	;

header:
	JANUS NUMBER_TOKEN ';'
	| JANUS NUMBER_TOKEN '{' header_options '}'
	;

header_options:
	/* empty */
	| header_options header_option
	;

header_option:
	expression_dot '=' expression ';'
	;


file_declaration:
	';'
	| import_statement
	| struct_declaration
	| interface_declaration
	| method_declaration
	| alias_declaration
	| operator_declaration
	| def_statement
	;

import_statement:
	IMPORT expression_dot ';'
	| IMPORT '.' '=' expression_dot ';'
	| IMPORT expression_dot '=' expression_dot ';'
	;

struct_declaration:
	STRUCT_OR_MSTRUCT type_name struct_options ';'
	| STRUCT_OR_MSTRUCT type_name struct_options '{' struct_content '}'
	;

struct_options:
	/* empty */
	| struct_options_
	;

struct_options_:
	struct_option
	| struct_options_ ',' struct_option
	;

struct_option:
	EXTENDS type
	| IMPLEMENTS type
	| IMPLEMENTS type ALIAS SYMBOL_TOKEN
	/* FIXME do we want   | SIZE expression */
	;

struct_content:
	/* empty */
	| struct_content struct_element
	;

struct_element:
	DEF SYMBOL_TOKEN type ';'
	| DEF SYMBOL_TOKEN function_type '{' function_content '}'
	| DEF SYMBOL_TOKEN type '{' function_content '}'
	| extends_declaration
	| implements_declaration
	;

extends_declaration:
	EXTENDS type ';'
	| EXTENDS type '{' extends_content '}'
	;

extends_content:
	/* empty */
	| extends_content extends_item
	;

extends_item:
	SYMBOL_TOKEN '=' SYMBOL_TOKEN ';'
	;

implements_declaration:
	IMPLEMENTS type ';'
	| IMPLEMENTS type '{' implements_content '}'
	| IMPLEMENTS type ALIAS SYMBOL_TOKEN ';'
	| IMPLEMENTS type ALIAS SYMBOL_TOKEN '{' implements_content '}'
	;

implements_content:
	/* empty */
	| implements_content implements_item
	;

implements_item:
	SYMBOL_TOKEN '=' SYMBOL_TOKEN ';'
	;

type_name:
	SYMBOL_TOKEN
	| SYMBOL_TOKEN '(' parameter_list ')'
	;

interface_declaration:
	INTERFACE type_name interface_options '{' interface_content '}'
	;

interface_options:
	/* empty */
	| interface_options_
	;

interface_options_:
	interface_option
	| interface_options_ ',' interface_option
	;

interface_option:
	EXTENDS type
	;

interface_content:
	/* empty */
	| interface_content interface_element
	;

interface_element:
	DEF SYMBOL_TOKEN type ';'
	| DEF SYMBOL_TOKEN function_type ';'
	| extends_declaration
	;

method_declaration:
	METHOD type SYMBOL_TOKEN function_type '{' function_content '}'
	| METHOD type SYMBOL_TOKEN function_type '=' expression ';'
	| METHOD type SYMBOL_TOKEN function_type ';'
	| METHOD type SYMBOL_TOKEN type '{' function_content '}'
	| METHOD type SYMBOL_TOKEN type '=' expression ';'
	| METHOD type SYMBOL_TOKEN type ';'
	;

alias_declaration:
	ALIAS SYMBOL_TOKEN type ';'
	;

operator_declaration:
	OPERATOR ANY_OP function_type '{' function_content '}'
	| OPERATOR ANY_OP function_type '=' expression ';'
	| OPERATOR ANY_OP function_type ';'
	| OPERATOR ANY_OP type '{' function_content '}'
	| OPERATOR ANY_OP type '=' expression ';'
	| OPERATOR ANY_OP type ';'
	;

def_statement:
	def_or_const SYMBOL_TOKEN function_type initializer
	| def_or_const SYMBOL_TOKEN function_type '{' function_content '}'
	| def_or_const SYMBOL_TOKEN function_type ';'
	| def_or_const SYMBOL_TOKEN type initializer
	| def_or_const SYMBOL_TOKEN type '{' function_content '}'
	| def_or_const SYMBOL_TOKEN type ';'
	| def_or_const SYMBOL_TOKEN initializer
	;

def_or_const:
	DEF
	| CONST
	;

initializer:
	'=' expression ';'
	| '=' '{' map_content '}'
	| '=' '[' list_content ']'
	;

function_content:
	/* empty */
	| function_content function_statement
	;

function_statement:
	';'
	| '{' function_content '}'
	| def_statement
	| IF expression '{' function_content '}'
	| IF expression '{' function_content '}' else_statement
	| WHILE expression '{' function_content '}'
	| FOR SYMBOL_TOKEN '=' expression '{' function_content '}'
	| WITH SYMBOL_TOKEN '=' expression '{' function_content '}'
	| RETURN expression ';'
	| RETURN ';'
	| CONTINUE ';'
	| BREAK NUMBER_TOKEN ';'
	| BREAK ';'
	| LABEL SYMBOL_TOKEN ';'
	| GOTO SYMBOL_TOKEN ';'
	| assignment_statement
	;

else_statement:
	ELSE '{' function_content '}'
	| ELSE IF expression '{' function_content '}'
	| ELSE IF expression '{' function_content '}' else_statement
	;

assignment_statement:
	expression ';'
	| expression initializer
	| expression ASSIGNMENT_OP expression ';'
	;

expression:
	expression_and
	| expression OR_OP expression_and   //  "|"  "^|"
	;

expression_and:
	expression_compare
	| expression_and AND_OP expression_compare  //  "&"
	;

expression_compare:
	expression_add
	| expression_compare COMPARE_OP expression_add //"==" "!=" "~~" "!~" "<=" ">=" ">" "<" ":"
	;

expression_add:
	expression_mult
	| expression_add ADD_OP expression_mult  //  "+" "-"
	;

expression_mult:
	expression_exp
	| expression_mult MULT_OP expression_exp // "*" "/"  "//" "+/" "-/" "%%" "+%" "-%" "<<" ">>"
	;

expression_exp:
	expression_prefix
	| expression_exp EXP_OP expression_prefix  //  "^"

expression_prefix:
	expression_suffix
	| PREFIX_OP expression_prefix   // "!"

expression_suffix:
	NUMBER_TOKEN
	| STRING_TOKEN
	| CHARACTER_TOKEN
	| FUNCTION function_type
	| expression_suffix SUFFIX_OP
	| expression_suffix '[' expression ']'
	| expression_suffix '(' list_content ')'
	| '(' expression ')'
	| expression_dot
	;

expression_dot:
	SYMBOL_TOKEN
	| expression_dot '.' SYMBOL_TOKEN
	;

type:
	expression_dot
	expression_dot '(' list_content ')'
	FUNCTION function_type
	;

function_type:
	'(' parameter_list ')'
	| '(' parameter_list ')' "->" type
	;

list_content:
	/* empty */
	| list_content_
	;

list_content_:
	expression
	| list_content_ ',' expression
	;

map_content:
	/* empty */
	| expression '=' expression ';' map_content

parameter_list:
	/* empty */
	| parameter_list_
	;

parameter_list_:
	SYMBOL_TOKEN type
	| SYMBOL_TOKEN '>' type
	| parameter_list_ ',' SYMBOL_TOKEN type
	| parameter_list_ ',' SYMBOL_TOKEN '>' type
	;

%%

