
#ifndef INCLUDED_PARSER_H
#define INCLUDED_PARSER_H

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include "compiler/parser_file.h"
#include "compiler/parser_node.h"

typedef struct ParserContext ParserContext; // FIXME maybe not needed?

void ParseSetDebug(bool on);

ParserNode *ParseFile(ParserFile *file, ParserContext *context);

// symbols created from tokens
extern ParserSymbol SYM_UNDEF;
extern ParserSymbol SYM_IDENTIFIER;
extern ParserSymbol SYM_NUMBER;
extern ParserSymbol SYM_CHARCONST;
extern ParserSymbol SYM_STRINGCONST;
extern ParserSymbol SYM_PUNCTUATION;
extern ParserSymbol SYM_KEYWORD;
extern ParserSymbol SYM_OPERATOR;

extern ParserSymbol SYM_FIXME;

// production symbols
extern ParserSymbol SYM_EMPTY; // for rules that don't have any content
extern ParserSymbol SYM_LIST;
extern ParserSymbol SYM_DOT_OP;
extern ParserSymbol SYM_IMPORT;
extern ParserSymbol SYM_IMPORT_PRIVATE;
extern ParserSymbol SYM_PROTOTYPE;
extern ParserSymbol SYM_FUNC;
extern ParserSymbol SYM_DECLARATION;
extern ParserSymbol SYM_USING;
extern ParserSymbol SYM_USING_AS;
extern ParserSymbol SYM_TRAILING_COMMA;
extern ParserSymbol SYM_ELLIPSIS;
extern ParserSymbol SYM_PARAMETER;
extern ParserSymbol SYM_PARAM_TYPE;
extern ParserSymbol SYM_DECL_TYPE;
extern ParserSymbol SYM_INITIALIZE;
extern ParserSymbol SYM_INIT_STRUCT;
extern ParserSymbol SYM_INIT_ARRAY;
extern ParserSymbol SYM_STATEMENT_LIST;
extern ParserSymbol SYM_STRUCT_LIST;
extern ParserSymbol SYM_ENUM_LIST;
extern ParserSymbol SYM_ENUM_ELEMENT;
extern ParserSymbol SYM_STRUCT_DEC;
extern ParserSymbol SYM_STRUCT_DEF;
extern ParserSymbol SYM_UNION_DEC;
extern ParserSymbol SYM_UNION_DEF;
extern ParserSymbol SYM_ENUM_DEC;
extern ParserSymbol SYM_ENUM_DEF;
extern ParserSymbol SYM_EXPRESSION_STATEMENT;
extern ParserSymbol SYM_EMPTY_STATEMENT;
extern ParserSymbol SYM_LABEL_STATEMENT;
extern ParserSymbol SYM_FOR_STATEMENT;
extern ParserSymbol SYM_WHILE_STATEMENT;
extern ParserSymbol SYM_DO_STATEMENT;
extern ParserSymbol SYM_IF_STATEMENT;
extern ParserSymbol SYM_IF_ELSE;
extern ParserSymbol SYM_SWITCH_STATEMENT;
extern ParserSymbol SYM_BREAK_STATEMENT;
extern ParserSymbol SYM_CONTINUE_STATEMENT;
extern ParserSymbol SYM_GOTO_STATEMENT;
extern ParserSymbol SYM_RETURN_STATEMENT;
extern ParserSymbol SYM_RETURN_VOID;
extern ParserSymbol SYM_CASE_ELEMENT;
extern ParserSymbol SYM_CASE_END_ELEMENT;
extern ParserSymbol SYM_DEFAULT_LABEL;
extern ParserSymbol SYM_CASE_LABEL;
extern ParserSymbol SYM_CONSTANT;
extern ParserSymbol SYM_ASSIGN_OP;
extern ParserSymbol SYM_ASSIGN_MULT_OP;
extern ParserSymbol SYM_ASSIGN_DIV_OP;
extern ParserSymbol SYM_ASSIGN_MOD_OP;
extern ParserSymbol SYM_ASSIGN_ADD_OP;
extern ParserSymbol SYM_ASSIGN_SUB_OP;
extern ParserSymbol SYM_ASSIGN_SHR_OP;
extern ParserSymbol SYM_ASSIGN_SHL_OP;
extern ParserSymbol SYM_ASSIGN_AND_OP;
extern ParserSymbol SYM_ASSIGN_OR_OP;
extern ParserSymbol SYM_ASSIGN_XOR_OP;
extern ParserSymbol SYM_CONDITIONAL;
extern ParserSymbol SYM_LOG_OR_OP;
extern ParserSymbol SYM_LOG_AND_OP;
extern ParserSymbol SYM_OR_OP;
extern ParserSymbol SYM_AND_OP;
extern ParserSymbol SYM_ADDR_OP;
extern ParserSymbol SYM_XOR_OP;
extern ParserSymbol SYM_EQUAL_OP;
extern ParserSymbol SYM_NEQUAL_OP;
extern ParserSymbol SYM_LESS_OP;
extern ParserSymbol SYM_GREATER_OP;
extern ParserSymbol SYM_LESSEQ_OP;
extern ParserSymbol SYM_GREATEREQ_OP;
extern ParserSymbol SYM_SHL_OP;
extern ParserSymbol SYM_SHR_OP;
extern ParserSymbol SYM_ADD_OP;
extern ParserSymbol SYM_SUB_OP;
extern ParserSymbol SYM_DIV_OP;
extern ParserSymbol SYM_MOD_OP;
extern ParserSymbol SYM_MULT_OP;
extern ParserSymbol SYM_PTR_OP;
extern ParserSymbol SYM_NOT_OP;
extern ParserSymbol SYM_BITNOT_OP;
extern ParserSymbol SYM_PREINC_OP;
extern ParserSymbol SYM_PREDEC_OP;
extern ParserSymbol SYM_POSTINC_OP;
extern ParserSymbol SYM_POSTDEC_OP;
extern ParserSymbol SYM_NEG_OP;
extern ParserSymbol SYM_POS_OP;
extern ParserSymbol SYM_SIZEOF_OP;
extern ParserSymbol SYM_ARRAY_OP;
extern ParserSymbol SYM_CALL_OP;
extern ParserSymbol SYM_INIT_OP;
extern ParserSymbol SYM_PAREN_EXPRESSION;

extern ParserSymbol SYM_STRING;
extern ParserSymbol SYM_TYPE_EXPRESSION;
extern ParserSymbol SYM_TYPE_ARRAY;
extern ParserSymbol SYM_TYPE_ARRAY_MATCH;
extern ParserSymbol SYM_TYPE_BITFIELD;
extern ParserSymbol SYM_TYPE_LINKAGE;
extern ParserSymbol SYM_TYPE_LINKNAME;

#endif

