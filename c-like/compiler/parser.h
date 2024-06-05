
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

struct ParserSymbol
{
	const char *rule_name;
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

// production symbols
extern ParserSymbol SYM_EMPTY; // for rules that don't have any content

extern ParserSymbol SYM_FIXME;

#endif

