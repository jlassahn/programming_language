
#include "compiler/exit_codes.h"
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
#include <string.h>
#include <stdbool.h>
#include <stdint.h>

extern int yydebug;

static CompileState compile_state;

bool AddBaseDir(CompileState *cs, const char *path)
{
	if (!IsValidPath(path))
	{
		Error(ERROR_FILE,
				"parameter '%s' is not a valid path.", path);
		return false;
	}

	StringBuffer *sb = StringBufferFromChars(path);

	// FIXME maybe an official function for character replace
	for (int i=0; i<sb->string.length; i++)
	{
		if ((sb->buffer[i] == '/') || (sb->buffer[i] == '\\'))
			sb->buffer[i] = PATH_SEPARATOR;
	}
	if (sb->buffer[sb->string.length-1] != PATH_SEPARATOR)
	{
		sb = StringBufferAppendChars(sb, PATH_SEPARATOR_STRING);
	}

	if (!DoesDirectoryExist(sb->buffer))
	{
		Warning(ERROR_FILE, "path '%s' does not exist.", sb->buffer);
		StringBufferFree(sb);
		return false;
	}

	StringBufferLock(sb);
	ListInsertLast(&compile_state.basedirs, sb);
	return true;
}

bool CheckForModuleStem(StringBuffer *dir, StringBuffer *stem)
{
	bool ret = false;

	printf("searching %s for %s\n", dir->buffer, stem->buffer);
	if (!DoesDirectoryExist(dir->buffer))
		return false;

	DirectorySearch *ds = DirectorySearchStart(dir->buffer);
	if (!ds)
		return false;

	while (true)
	{
		const char *name = DirectorySearchNextFile(ds);
		if (name == NULL)
			break;

		if (strncmp(name, stem->buffer, stem->string.length) == 0)
		{
			printf("FOUND: %s\n", name);
			ret = true;
			break;
		}
	}
	DirectorySearchEnd(ds);

	return ret;
}

bool CheckForModuleFiles(List *base_paths, const char *name)
{
	bool ret = false;

	printf("FIXME checking module '%s' ...\n", name);

	int path_end = 0;
	StringBuffer *mod_path = StringBufferFromChars(name);

	// FIXME maybe an official function for character replace
	for (int i=0; i<mod_path->string.length; i++)
	{
		if (mod_path->buffer[i] == '.')
		{
			path_end = i+1;
			mod_path->buffer[i] = PATH_SEPARATOR;
		}
	}
	StringBuffer *stem = StringBufferFromChars(&mod_path->buffer[path_end]);
	stem = StringBufferAppendChars(stem, ".");

	// FIXME maybe an official function to trim a string buffer to length?
	mod_path->buffer[path_end] = 0;
	mod_path->string.length = path_end;

	StringBuffer *dir = StringBufferCreateEmpty(200);

	for (ListEntry *entry=base_paths->first; entry!=NULL; entry=entry->next)
	{
		StringBuffer *base = entry->item;

		StringBufferClear(dir);
		dir = StringBufferAppendBuffer(dir, base);
		dir = StringBufferAppendChars(dir, "source/");
		dir = StringBufferAppendBuffer(dir, mod_path);
		if (CheckForModuleStem(dir, stem))
		{
			ret = true;
			break;
		}

		StringBufferClear(dir);
		dir = StringBufferAppendBuffer(dir, base);
		dir = StringBufferAppendChars(dir, "import/");
		dir = StringBufferAppendBuffer(dir, mod_path);
		if (CheckForModuleStem(dir, stem))
		{
			ret = true;
			break;
		}
	}

	StringBufferFree(dir);
	StringBufferFree(stem);
	StringBufferFree(mod_path);
	return ret;
}


bool AddInputModule(CompileState *cs, const char *name)
{
	if (!IsValidNamespace(name))
		return false;

	if (!CheckForModuleFiles(&cs->basedirs, name))
		return false;

	StringBuffer *sb = StringBufferFromChars(name);
	StringBufferLock(sb);
	ListInsertLast(&compile_state.input_modules, sb);

	return true;
}

bool AddInputFile(CompileState *cs, const char *name)
{
	if (!DoesFileExist(name))
		return false;

	StringBuffer *sb = StringBufferFromChars(name);
	StringBufferLock(sb);

	// FIXME if name ends with ".moss" assume it's a source file
	//       otherwise it's a library or something.

	ListInsertLast(&compile_state.input_files, sb);

	return true;
}

bool AddInput(CompileState *cs, const char *name)
{
	if (!IsValidPath(name))
	{
		Error(ERROR_FILE,
				"parameter '%s' is not a valid input name.", name);
		return false;
	}

	if (AddInputModule(cs, name))
		return true;
	if (AddInputFile(cs, name))
		return true;

	Error(ERROR_FILE,
			"Input name '%s' not found as either a file or a module.", name);
	return false;
}

int main(int argc, const char *argv[])
{
	const CompilerArgs *args = ParseArgs(argc, argv);
	if (args == NULL)
		return EXIT_USAGE;

	// check args for validity
	//    warnings;
	//    optimizations;
	//    generation;
	//    defines;
	//    versions;
	//    outfile;
	//    outdir;
	//    treefile;

	// FIXME also add "well known" search paths here.  How? Environment vars?
	AddBaseDir(&compile_state, ".");
	for (ArgStringList *entry=args->basedirs; entry!=NULL; entry=entry->next)
		AddBaseDir(&compile_state, entry->arg);

	// must be done after basedirs are finished, so it can search for modules
	bool inputs_good = true;
	for (ArgStringList *entry=args->inputs; entry!=NULL; entry=entry->next)
	{
		inputs_good = AddInput(&compile_state, entry->arg) && inputs_good;
	}


	// FIXME printing results...
	if (!inputs_good)
	{
		printf("BAD INPUTS\n");
		// FIXME skip compile steps and exit
	}

	for (ListEntry *entry=compile_state.basedirs.first;
			entry!=NULL; entry=entry->next)
	{
		StringBuffer *sb = entry->item;
		printf("search directory: %s\n", sb->string.data);
	}

	for (ListEntry *entry=compile_state.input_files.first;
			entry!=NULL; entry=entry->next)
	{
		StringBuffer *sb = entry->item;
		printf("input file: %s\n", sb->string.data);
	}

	for (ListEntry *entry=compile_state.input_modules.first;
			entry!=NULL; entry=entry->next)
	{
		StringBuffer *sb = entry->item;
		printf("input module: %s\n", sb->string.data);
	}


	// FIXME cleanup starts here

	CompileStateFree(&compile_state);
	FreeArgs(args);

	printf("allocation count = %d\n", AllocCount());

	printf("errors: %d, warnings: %d\n", ErrorCount(), WarningCount());

	if (ErrorCount() > 0)
		return EXIT_DATAERR;
	return EXIT_OK;

	//yydebug = 1;
	const char * filename = "examples/source/hello.moss";
	if (argc == 2)
		filename = argv[1];

	ParserFile *file = FileRead(filename);
	if (!file)
	{
		printf("can't open file %s\n", filename);
		return EXIT_NOINPUT;
	}

	ParserNode *root = ParseFile(file, NULL);
	PrintNodeTree(stdout, root);
	printf("nodes = %d\n", GetNodeCount());
	FreeNode(root);
	printf("nodes = %d\n", GetNodeCount());

	FileFree(file);

	if (ErrorCount() > 0)
		return EXIT_DATAERR;
	return EXIT_OK;
}

