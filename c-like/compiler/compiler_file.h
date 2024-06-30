
#ifndef INCLUDED_COMPILER_FILE_H
#define INCLUDED_COMPILER_FILE_H

#include "compiler/parser_file.h"
#include "compiler/parser.h"
#include "compiler/types.h"
#include "compiler/namespace.h"

typedef struct CompilerFile CompilerFile;
struct CompilerFile
{
	uint32_t flags;
	StringBuffer *path;
	ParserFile *parser_file;
	ParserNode *root;
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

