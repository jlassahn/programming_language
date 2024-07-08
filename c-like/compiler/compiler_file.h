
#ifndef INCLUDED_COMPILER_FILE_H
#define INCLUDED_COMPILER_FILE_H

#include "compiler/parser_file.h"
#include "compiler/parser.h"
#include "compiler/types.h"
#include "compiler/namespace.h"

typedef struct ImportLink ImportLink;
struct ImportLink
{
	ParserNode *parse;
	bool is_private;
	Namespace *namespace;
};

typedef struct CompilerFile CompilerFile;
struct CompilerFile
{
	ParserFile parser_file;
	uint32_t flags;
	StringBuffer *path; // FIXME duplicated info in parser_file
	ParserNode *root;
	Namespace *namespace;

	List imports; // List of ImportLink*
};


typedef enum
{
	FILE_FROM_INPUT = 0x0001,
	FILE_PARSE_FAILED = 0x0002,
}
CompilerFileFlags;

CompilerFile *CompilerFileCreate(StringBuffer *path);
void CompilerFileFree(CompilerFile *cf);

bool CompilerFilePickNamespace(CompilerFile *cf, Namespace *root);

#endif

