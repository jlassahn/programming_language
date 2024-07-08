
#include "compiler/pass_search_and_parse.h"
#include "compiler/compiler_file.h"
#include "compiler/parser.h"
#include "compiler/errors.h"
#include "compiler/search.h"
#include "compiler/memory.h"

const char *moss_file_extensions[] =
{
	".moss",
	NULL
};

bool ScanFileImports(CompilerFile *cf, CompileState *state)
{
	printf("FIXME scanning file imports for %s\n", cf->path->buffer);
	return true; // FIXME fake
}

void TranslateFileScopeItem(ParserNode *node, CompilerFile *cf,
		CompileState *state)
{
	if (node->symbol == &SYM_IMPORT)
	{
		ImportLink *import = Alloc(sizeof(ImportLink));
		import->parse = node;
		import->is_private = false;
		import->namespace = NULL; // FIXME find
		ListInsertLast(&cf->imports, import);
	}
	else if (node->symbol == &SYM_IMPORT_PRIVATE)
	{
		ImportLink *import = Alloc(sizeof(ImportLink));
		import->parse = node;
		import->is_private = true;
		import->namespace = NULL; // FIXME find
		ListInsertLast(&cf->imports, import);
	}
	else
	{
		printf("FIXME skipping top-level node %s\n", node->symbol->rule_name);
		return;
	}
}

void TranslateFileScopeList(ParserNode *node, CompilerFile *cf,
		CompileState *state)
{
	if ((node->symbol == &SYM_LIST) && (node->count == 2))
	{
		TranslateFileScopeList(node->children[0], cf, state);
		ParserNode *item = node->children[1];
		TranslateFileScopeItem(item, cf, state);
	}
	else if (node->symbol != &SYM_EMPTY)
	{
		ParserNode *item = node;
		TranslateFileScopeItem(item, cf, state);
	}
}

void TopLevelScan(CompilerFile *cf, CompileState *state)
{
	TranslateFileScopeList(cf->root, cf, state);
}

bool ParseInputFile(CompilerFile *cf, CompileState *state)
{
	Namespace *root = &state->root_namespace;

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

	if (cf->root != NULL)
		TopLevelScan(cf, state);
	if (!ScanFileImports(cf, state))
		return false;

	return true;
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

	if (cf->root != NULL)
		TopLevelScan(cf, state);
	if (!ScanFileImports(cf, state))
		ret = false;

	return ret;
}


bool ScanNamespaceFiles(Namespace *ns, CompileState *state)
{
	if (ns->flags & NAMESPACE_SCANNED)
		return true;
	ns->flags |= NAMESPACE_SCANNED;

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
		if (!DoModuleFile(file, ns, is_private, state))
			ret = false;
	}
	SearchFilesEnd(sf);

	StringBufferFree(stem);

	if ((ns->private_files.first == NULL) && (ns->public_files.first == NULL))
	{
		Error(ERROR_FILE, "No files found for module %s\n", ns->path->buffer);
		ret = false;
	}

	return ret;
}

bool PassSearchAndParse(CompileState *state)
{
	bool ret = true;

	for (ListEntry *entry=state->input_files.first;
			entry!=NULL; entry=entry->next)
	{
		CompilerFile *cf = entry->item;

		if (!ParseInputFile(cf, state))
		{
			ret = false;
			continue;
		}
		// FIXME scan namespace for other files
	}

	for (ListEntry *entry=state->input_modules.first;
			entry!=NULL; entry=entry->next)
	{
		Namespace *module = entry->item;
		if (!ScanNamespaceFiles(module, state))
			ret = false;
	}

	return ret;
}

