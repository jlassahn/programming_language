
#include "parser_file.h"
#include "tokenizer.h"
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

int main(int argc, const char *argv[])
{
	const char * filename = "examples/source/hello.moss";
	if (argc == 2)
		filename = argv[1];

	ParserFile *file = FileRead(filename);
	if (!file)
	{
		printf("can't open file %s\n", filename);
		return 1;
	}

	Tokenizer tokenizer;
	Token token;

	TokenizerStart(&tokenizer, file);
	while (!TokenizerIsEOF(&tokenizer))
	{
		GetCurrentToken(&tokenizer, &token);
		TokenizerConsume(&tokenizer);

		TokenPrint(stdout, &token);
	}

	FileFree(file);

	return 0;
}

