
#include "compiler/exit_codes.h"
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/fileio.h"
#include "compiler/errors.h"
#include "compiler/commandargs.h"
#include "compiler/parser_file.h"
#include "compiler/compiler_file.h"
#include "compiler/tokenizer.h"
#include "compiler/parser.h"
#include "compiler/compile_state.h"
#include "compiler/namespace.h"
#include "compiler/search.h"
#include "compiler/stringtypes.h"
#include "compiler/pass_configure.h"
#include "compiler/pass_search_and_parse.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>

static CompileState compile_state;

static const char *moss_file_extensions[] =
{
	".moss",
	NULL
};


bool ParseInputFile(CompilerFile *cf, Namespace *root)
{
	if (!ParserFileRead(&cf->parser_file, cf->path->buffer))
		return false;

	cf->root = ParseFile(&cf->parser_file, NULL);
	if (cf->parser_file.parser_result != 0)
		cf->flags |= FILE_PARSE_FAILED;

	// determine namespace after parsing, in case we add a file
	// header that overrides the default filename-based namespace.
	if (!CompilerFilePickNamespace(cf, root))
	{
		Error(ERROR_FILE,
			"File name '%s' isn't a valid namespace.", cf->path->buffer);
		return false;
	}
	return true;
}

// FIXME needs namespace of file being scanned
bool ScanImportNodes(ParserNode *node, CompileState *state)
{
	// FIXME rename, reorganize, create List of ImportLink
	//  on files.
	if (node == NULL)
		return true;

	if (node->symbol == &SYM_EMPTY)
		return true;

	if (node->symbol == &SYM_LIST)
	{
		for (int i=0; i<node->count; i++)
		{
			ScanImportNodes(node->children[i], state);
		}
	}

	if ((node->symbol == &SYM_IMPORT) || (node->symbol == &SYM_IMPORT_PRIVATE))
	{
		int start = node->position.start.offset;
		int end = node->position.end.offset;
		const char *data = node->position.file->data;
		printf("FIXME scanning import node %s: %.*s\n", node->symbol->rule_name, end-start, data+start);
	}

	return true;
}

bool ScanFileImports(CompilerFile *cf, CompileState *state)
{
	// FIXME mark file as done with imports
	// FIXME for each import ScanNamespaceFiles()
	printf("FIXME scanning %s for imports\n", cf->path->buffer);

	return ScanImportNodes(cf->root, state);
}

bool DoModuleFile(StringBuffer *path, Namespace *ns, bool is_private,
		CompileState *state)
{
	StringBufferLock(path);
	CompilerFile *cf = CompilerFileCreate(path);

	if (!ParserFileRead(&cf->parser_file, cf->path->buffer))
	{
		CompilerFileFree(cf);
		return false;
	}

	bool ret = true;
	cf->root = ParseFile(&cf->parser_file, NULL);
	if (cf->parser_file.parser_result != 0)
	{
		cf->flags |= FILE_PARSE_FAILED;
		ret = false;
	}

	cf->namespace = ns;
	if (is_private)
		ListInsertLast(&ns->private_files, cf);
	else
		ListInsertLast(&ns->public_files, cf);

	if (!ScanFileImports(cf, state))
		ret = false;

	return ret;
}

bool ScanNamespaceFiles(Namespace *ns, CompileState *state)
{
	// FIXME mark namespace as already scanned
	printf("FIXME scanning namespace %s (%s, %.*s)\n", ns->path->buffer, ns->parent->path->buffer, ns->stem.length, ns->stem.data);

	bool ret = true;
	List *base_paths = &state->basedirs;

	StringBuffer *stem = StringBufferFromString(&ns->stem);
	stem = StringBufferAppendChars(stem, ".");

	SearchFiles *sf;
	StringBuffer *file;

	sf = SearchFilesStart(base_paths, "source/",
			ns->parent->path->buffer, stem->buffer, moss_file_extensions);
	while (true)
	{
		bool is_private = true;
		file = SearchFilesNext(sf);
		if (file == NULL)
			break;
		printf("   found %s\n", file->buffer);
		if (!DoModuleFile(file, ns, is_private, state))
			ret = false;
	}
	SearchFilesEnd(sf);

	sf = SearchFilesStart(base_paths, "import/",
			ns->parent->path->buffer, stem->buffer, moss_file_extensions);
	while (true)
	{
		bool is_private = false;
		file = SearchFilesNext(sf);
		if (file == NULL)
			break;
		printf("   found %s\n", file->buffer);
		if (!DoModuleFile(file, ns, is_private, state))
			ret = false;
	}
	SearchFilesEnd(sf);

	StringBufferFree(stem);
	return ret;
}

int main(int argc, const char *argv[])
{
	// FIXME naming things ...
	// Scan...  Parse...
	// Translate...
	// Generate

	CompileStateInit(&compile_state);

	const CompilerArgs *args = ParseArgs(argc, argv);
	if (args == NULL)
		return EXIT_USAGE;

	const char *env = getenv("MOSS_IMPORT_PATH");
	if (env == NULL)
		env = "";

	if (!PassConfigure(&compile_state, args, env))
		return EXIT_USAGE;

	CompileStatePrint(&compile_state);

	ParseSetDebug(false);
	bool inputs_good = true;
	//bool inputs_good = PassSearchAndParse(&compile_state);

	// FIXME PassSearchAndParse(&compile_state);
	// FIXME rename things Parse... and Scan...
	for (ListEntry *entry=compile_state.input_files.first;
			entry!=NULL; entry=entry->next)
	{
		CompilerFile *cf = entry->item;

		if (!ParseInputFile(cf, &compile_state.root_namespace))
			inputs_good = false;

		if (!ScanFileImports(cf, &compile_state))
			inputs_good = false;

		// FIXME scan namespace for other files
	}

	for (ListEntry *entry=compile_state.input_modules.first;
			entry!=NULL; entry=entry->next)
	{
		Namespace *module = entry->item;
		if (!ScanNamespaceFiles(module, &compile_state))
			inputs_good = false;
	}

	// FIXME all file input is done.
	if (!inputs_good)
	{
		printf("BAD INPUTS\n");
		// FIXME skip compile steps and exit
	}

	// FIXME PassTranslate(&compile_state);
	// FIXME PassGenerate(&compile_state);
	// FIXME PassLink(&compile_state);

	int depth = 1;
	printf("\nNAMESPACE:\n");
	MapIterate(&compile_state.root_namespace.children,
			NamespacePrinter, &depth);

	// FIXME cleanup starts here

	CompileStateFree(&compile_state);
	FreeArgs(args);

	printf("allocation count = %d\n", AllocCount());

	printf("errors: %d, warnings: %d\n", ErrorCount(), WarningCount());

	if (ErrorCount() > 0)
		return EXIT_DATAERR;
	return EXIT_OK;
}

