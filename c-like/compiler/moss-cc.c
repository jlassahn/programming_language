
#include "compiler/parser_file.h"
#include "compiler/tokenizer.h"
#include "compiler/parser.h"
#include <stdio.h>

typedef struct StringList StringList;
typedef struct CompilerSettings CompilerSettings;

struct StringList
{
	const char **list;
	int count;
};

struct CompilerSettings
{
	StringList import_paths;
	StringList source_paths;
	StringList lib_paths;
	StringList targets;
};

extern int yydebug;

int main(int argc, const char *argv[])
{
	//yydebug = 1;
	const char * filename = "examples/source/hello.moss";
	if (argc == 2)
		filename = argv[1];

	ParserFile *file = FileRead(filename);
	if (!file)
	{
		printf("can't open file %s\n", filename);
		return 1;
	}

	ParserNode *root = ParseFile(file, NULL);
	PrintNodeTree(stdout, root);

	FileFree(file);

	return 0;
}

