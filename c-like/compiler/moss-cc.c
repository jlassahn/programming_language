
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
#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>

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
	ListInsertLast(&cs->basedirs, sb);
	return true;
}

bool CheckForModuleStem(StringBuffer *dir, StringBuffer *stem)
{
	bool ret = false;

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

bool AddInputModule(CompileState *state, const char *path)
{
	if (!IsValidNamespace(path))
		return false;

	if (!CheckForModuleFiles(&state->basedirs, path))
		return false;

	Namespace *ns = &state->root_namespace;

	String name;
	int length = strlen(path);
	int start = 0;
	int end = 0;
	for (int i=0; i<length; i++)
	{
		if (path[i] == '.')
		{
			end = i;
			name.data = &path[start];
			name.length = end-start;
			ns = NamespaceGetChild(ns, &name);
			if (ns == NULL)
				return false;
			start = i+1;
		}
	}
	if (start < length)
	{
		end = length;
		name.data = &path[start];
		name.length = end-start;
		ns = NamespaceGetChild(ns, &name);
		if (ns == NULL)
			return false;
	}

	ListInsertLast(&compile_state.input_modules, ns);

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

	CompilerFile *cf = CompilerFileCreate(sb);
	cf->flags |= FILE_FROM_INPUT;
	ListInsertLast(&compile_state.input_files, cf);

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

bool ParseInputFile(CompilerFile *cf, Namespace *root)
{
	cf->parser_file = FileRead(cf->path->buffer);
	if (!cf->parser_file)
		return false;

	cf->root = ParseFile(cf->parser_file, NULL);
	if (cf->parser_file->parser_result != 0)
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

bool ScanImportNodes(ParserNode *node, CompileState *state)
{
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
	// FIXME for each import ScanNamespaceFiles()
	printf("FIXME scanning %s for imports\n", cf->path->buffer);

	return ScanImportNodes(cf->root, state);
}

bool ScanModuleFiles(StringBuffer *dir, StringBuffer *stem,
		Namespace *ns, bool is_private, CompileState *state)
{
	printf("FIXME scanning module files %s, %s\n", dir->buffer, stem->buffer);

	bool ret = true;

	DirectorySearch *ds = DirectorySearchStart(dir->buffer);
	if (!ds)
		return true;

	while (true)
	{
		const char *name = DirectorySearchNextFile(ds);
		if (name == NULL)
			break;

		if (strncmp(name, stem->buffer, stem->string.length) == 0)
		{
			printf("   found %s\n", name);
			StringBuffer *sb = StringBufferFromString(&dir->string);
			sb = StringBufferAppendChars(sb, name);
			StringBufferLock(sb);
			CompilerFile *cf = CompilerFileCreate(sb);

			cf->parser_file = FileRead(cf->path->buffer);
			if (!cf->parser_file)
			{
				ret = false;
				CompilerFileFree(cf);
				continue;
			}

			cf->root = ParseFile(cf->parser_file, NULL);
			if (cf->parser_file->parser_result != 0)
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
		}
	}
	DirectorySearchEnd(ds);
	return ret;
}

bool ScanNamespaceFiles(Namespace *ns, CompileState *state)
{
	printf("FIXME scanning namespace %s (%s, %.*s)\n", ns->path->buffer, ns->parent->path->buffer, ns->stem.length, ns->stem.data);

	bool ret = true;
	List *base_paths = &state->basedirs;

	StringBuffer *dir = StringBufferCreateEmpty(200);
	StringBuffer *stem = StringBufferFromString(&ns->stem);
	stem = StringBufferAppendChars(stem, ".");
	for (ListEntry *entry=base_paths->first; entry!=NULL; entry=entry->next)
	{
		StringBuffer *base = entry->item;

		StringBufferClear(dir);
		dir = StringBufferAppendBuffer(dir, base);
		dir = StringBufferAppendChars(dir, "source/");
		dir = StringBufferAppendBuffer(dir, ns->parent->path);
		if (!ScanModuleFiles(dir, stem, ns, true, state))
			ret = false;

		StringBufferClear(dir);
		dir = StringBufferAppendBuffer(dir, base);
		dir = StringBufferAppendChars(dir, "import/");
		dir = StringBufferAppendBuffer(dir, ns->parent->path);
		if (!ScanModuleFiles(dir, stem, ns, false, state))
			ret = false;
	}

	StringBufferFree(stem);
	StringBufferFree(dir);
	return ret;
}

int main(int argc, const char *argv[])
{
	const CompilerArgs *args = ParseArgs(argc, argv);
	if (args == NULL)
		return EXIT_USAGE;

	ParseSetDebug(false);

	CompileStateInit(&compile_state);

	// FIXME handle other args
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

	CompileStatePrint(&compile_state);

	for (ListEntry *entry=compile_state.input_files.first;
			entry!=NULL; entry=entry->next)
	{
		CompilerFile *cf = entry->item;

		if (!ParseInputFile(cf, &compile_state.root_namespace))
			inputs_good = false;

		if (!ScanFileImports(cf, &compile_state))
			inputs_good = false;
	}

	for (ListEntry *entry=compile_state.input_modules.first;
			entry!=NULL; entry=entry->next)
	{
		Namespace *module = entry->item;
		if (!ScanNamespaceFiles(module, &compile_state))
			inputs_good = false;
	}

	if (!inputs_good)
	{
		printf("BAD INPUTS\n");
		// FIXME skip compile steps and exit
	}

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

