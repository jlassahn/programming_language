
#ifndef INCLUDED_PARSER_H
#define INCLUDED_PARSER_H

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include "compiler/parser_file.h"

typedef struct ParserSymbol ParserSymbol;
typedef struct ParserNode ParserNode;
typedef struct ParserContext ParserContext;

#define MAX_CHILDREN 10

enum ParseSymbolFlags
{
	PRINT_CONTENT = 0x0001,
};

struct ParserSymbol
{
	const char *rule_name;
	uint32_t flags;
	// FIXME include actions for generators here.
};

struct ParserNode
{
	ParserSymbol *symbol;
	FilePositionRange position;
	ParserNode *children[MAX_CHILDREN];
	int count;
	uint32_t flags; // e.g. constant value

	// extra data can include
	// number or sting values for constants
	// symbol table entries for variables
	// operator type for unary or binary operators
	// error information
};

ParserNode *ParseFile(ParserFile *file, ParserContext *context);

ParserNode *MakeNode(ParserSymbol *kind, int count, ParserNode **params);
void PrintNodeTree(FILE *fp, ParserNode *root);

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

#endif

