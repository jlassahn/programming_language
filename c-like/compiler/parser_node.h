
#ifndef INCLUDED_PARSER_NODE_H
#define INCLUDED_PARSER_NODE_H

#include "compiler/parser_file.h"

#include <stdio.h>

#define MAX_CHILDREN 10

enum ParseSymbolFlags
{
	PRINT_CONTENT = 0x0001,
	SYM_DISCARD = 0x0100,
};

typedef struct ParserSymbol ParserSymbol;
struct ParserSymbol
{
	const char *rule_name;
	uint32_t flags;
	// FIXME include actions for generators here.
};

typedef struct ParserNode ParserNode;
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

ParserNode *MakeNode(ParserSymbol *kind, int count, ParserNode **params);
void FreeNode(ParserNode *node);

void PrintNodeTree(FILE *fp, ParserNode *root);
int GetNodeCount(void);

ParserNode *GetLastNode(void);

#endif

