
#ifndef INCLUDED_PARSER_H
#define INCLUDED_PARSER_H

#include <stdint.h>
#include <stdbool.h>

typedef struct ParserSymbol ParserSymbol;
typedef struct ParserNode ParserNode;
typedef struct ParserContext ParserContext;

#define MAX_CHILDREN 8

struct ParserSymbol
{
	const char *rule_name;
};

struct ParserNode
{
	ParserSymbol *symbol;
	FilePositionRange position;
	ParserNode *children[MAX_CHILDREN];
	uint32_t flags; // e.g. constant value

	// extra data can include
	// number or sting values for constants
	// symbol table entries for variables
	// operator type for unary or binary operators
	// error information
};

ParserNode *ParseFile(ParserContext *context);

extern ParserSymbol *SYM_FILE;

#endif

