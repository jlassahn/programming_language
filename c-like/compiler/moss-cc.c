
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/fileio.h"
#include "compiler/errors.h"
#include "compiler/commandargs.h"
#include "compiler/parser_file.h"
#include "compiler/tokenizer.h"
#include "compiler/parser.h"
#include "compiler/compile_state.h"
#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>

static CompileState compile_state;


extern int yydebug;

int main(int argc, const char *argv[])
{
	const CompilerArgs *args = ParseArgs(argc, argv);
	if (args == NULL)
		return -1;

	// check args for validity
	//    warnings;
	//    optimizations;
	//    generation;
	//    defines;
	//    versions;
	//    outfile;
	//    outdir;
	//    treefile;

	// basedirs;
	ListInsertLast(&compile_state.basedirs, StringBufferFromChars("."));
	for (ArgStringList *entry=args->basedirs; entry!=NULL; entry=entry->next)
	{
		printf("basedir = %s\n", entry->arg);
		if (!IsValidPath(entry->arg))
		{
			Error(ERROR_FILE,
					"parameter '%s' is not a valid path.", entry->arg);
			continue;
		}
		if (!DoesDirectoryExist(entry->arg))
		{
			Error(ERROR_FILE, "path '%s' does not exist.", entry->arg);
			continue;
		}
		ListInsertLast(&compile_state.basedirs, StringBufferFromChars(entry->arg));
	}

	// inputs;

	printf("inputs:\n");
	PrintArgList(args->inputs);
	printf("defines:\n");
	PrintArgList(args->defines);
	printf("basedirs:\n");
	PrintArgList(args->basedirs);

	CompileStateFree(&compile_state);
	FreeArgs(args);

	printf("allocation count = %d\n", AllocCount());
	return 0;

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
	printf("nodes = %d\n", GetNodeCount());
	FreeNode(root);
	printf("nodes = %d\n", GetNodeCount());

	FileFree(file);

	return 0;
}

