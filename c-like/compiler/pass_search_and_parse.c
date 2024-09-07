
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

bool ScanNamespaceFiles(Namespace *ns, CompileState *state);

bool ScanFileImports(CompilerFile *cf, CompileState *state)
{
	bool ret = true;
	for (ListEntry *entry=cf->imports.first; entry!=NULL; entry=entry->next)
	{
		ImportLink *import = entry->item;
		if (!ScanNamespaceFiles(import->namespace, state))
			ret = false;
	}

	return ret;
}

Namespace *FindNamespaceFromDotList(ParserNode *node, Namespace *cns,
		Namespace *root_ns)
{
	String next_name;
	Namespace *cur;
	if ((node->symbol == &SYM_DOT_OP) && (node->count == 2))
	{
		cur = FindNamespaceFromDotList(node->children[0], cns, root_ns);
		if (cur == NULL)
			return NULL;

		if (!ParserNodeGetValue(node->children[1], &next_name))
			return NULL;

		// special handling for "_"
		if (StringEqualsCString(&next_name, "_"))
			return cur->parent; //can be NULL for malformed import lists

		return NamespaceMakeChild(cur, &next_name);
	}

	// root node
	if (!ParserNodeGetValue(node, &next_name))
		return NULL;

	// special handling for "_" at beginning of list
	if (StringEqualsCString(&next_name, "_"))
		return cns->parent;

	// FIXME check for identifiers that shadow predefined names
	return NamespaceMakeChild(root_ns, &next_name);
}

ImportLink *TranslateImportLink(ParserNode *node, Namespace *cns,
		Namespace *root_ns)
{
	if (node->count != 1)
		return NULL;

	Namespace *ns = FindNamespaceFromDotList(node->children[0], cns, root_ns);
	if ((ns == NULL) || (ns == root_ns))
	{
		ErrorAt(ERROR_FILE, node->position.file->filename, &node->position.start,
				"import with invalid namespace name");
		return NULL;
	}

	ImportLink *import = Alloc(sizeof(ImportLink));
	import->parse = node;
	import->is_private = false;
	import->namespace = ns;

	return import;
}

bool TranslateDeclaration(ParserNode *node, Namespace *ns, bool is_private)
{
	ParserNode *dtype = node->children[0];
	ParserNode *name = node->children[1];
	ParserNode *properties = node->children[2];
	ParserNode *value = node->children[3];

	String name_str;
	if (!ParserNodeGetValue(name, &name_str))
	{
		// FIXME can't happen?
		return false;
	}

	Symbol *sym_pub = MapFind(&ns->public_syms.exports, &name_str);
	Symbol *sym_all = MapFind(&ns->private_syms.exports, &name_str);

	if (sym_all == NULL)
	{
		sym_all = SymbolCreate(&name_str);
		MapInsert(&ns->private_syms.exports, &name_str, sym_all);
	}

	if (!is_private && (sym_pub == NULL))
	{
		sym_pub = SymbolCreate(&name_str);
		sym_pub->associated = sym_all;
		sym_all->associated = sym_pub;
		MapInsert(&ns->public_syms.exports, &name_str, sym_pub);
	}

	// FIXME store or merge more symbol information
	(void)dtype;
	(void)properties;
	(void)value;

	return true;
}

void TranslateFileScopeItem(ParserNode *node, CompilerFile *cf, bool is_private,
		CompileState *state)
{
	if (node->symbol == &SYM_IMPORT)
	{
		ImportLink *import = TranslateImportLink(node, cf->namespace,
				&state->root_namespace);
		if (import != NULL)
		{
			import->is_private = false;
			ListInsertLast(&cf->imports, import);
			ListInsertLast(&cf->namespace->private_syms.imports, import);
			if (!is_private)
				ListInsertLast(&cf->namespace->public_syms.imports, import);
		}
	}
	else if (node->symbol == &SYM_IMPORT_PRIVATE)
	{
		ImportLink *import = TranslateImportLink(node, cf->namespace,
				&state->root_namespace);
		if (import != NULL)
		{
			import->is_private = true;
			ListInsertLast(&cf->imports, import);
			ListInsertLast(&cf->namespace->private_syms.imports, import);
			if (!is_private)
				ListInsertLast(&cf->namespace->public_syms.imports, import);
		}
	}
	else if (node->symbol == &SYM_DECLARATION)
	{
		TranslateDeclaration(node, cf->namespace, is_private);
	}
	else
	{
		// FIXME handle other top-level stuff
		// SYM_USING
		// SYM_USING_AS
		// SYM_PROTOTYPE
		// SYM_FUNC
		// SYM_STRUCT_DEC
		// SYM_STRUCT_DEF
		// SYM_UNION_DEC
		// SYM_UNION_DEF
		// SYM_ENUM_DEC
		// SYM_ENUM_DEF
		return;
	}
}

void TranslateFileScopeList(ParserNode *node, CompilerFile *cf, bool is_private,
		CompileState *state)
{
	if ((node->symbol == &SYM_LIST) && (node->count == 2))
	{
		TranslateFileScopeList(node->children[0], cf, is_private, state);
		ParserNode *item = node->children[1];
		TranslateFileScopeItem(item, cf, is_private, state);
	}
	else if (node->symbol != &SYM_EMPTY)
	{
		ParserNode *item = node;
		TranslateFileScopeItem(item, cf, is_private, state);
	}
}

void TopLevelScan(CompilerFile *cf, bool is_private, CompileState *state)
{
	TranslateFileScopeList(cf->root, cf, is_private, state);
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

	bool is_private = true; // explicitly specified input files are private
	if (cf->root != NULL)
		TopLevelScan(cf, is_private, state);
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
		ListInsertLast(&ns->private_syms.files, cf);
	else
		ListInsertLast(&ns->public_syms.files, cf);

	if (cf->root != NULL)
		TopLevelScan(cf, is_private, state);
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

	if ((ns->private_syms.files.first == NULL) && (ns->public_syms.files.first == NULL))
	{
		Error(ERROR_FILE, "No files found for module '%s'.", ns->path->buffer);
		ret = false;
	}

	return ret;
}

typedef struct AddGlobalsCtx AddGlobalsCtx;
struct AddGlobalsCtx
{
	bool ret;
	CompileState *state;
};

void AddGlobalsToNamespace(const String *key, void *value, void *vctx)
{
	Namespace *ns = value;
	AddGlobalsCtx *ctx = vctx;

	MapIterate(&ns->children, AddGlobalsToNamespace, vctx);

	bool ret = true;
	if (!SymbolTableInsertMap(&ns->private_syms.symbol_table, &ctx->state->builtins))
		ret = false;
	if (!SymbolTableInsertMap(&ns->private_syms.symbol_table, &ns->private_syms.exports))
		ret = false;

	if (!SymbolTableInsertMap(&ns->public_syms.symbol_table, &ctx->state->builtins))
		ret = false;
	if (!SymbolTableInsertMap(&ns->public_syms.symbol_table, &ns->public_syms.exports))
		ret = false;

	if (!ret)
		ctx->ret = false;
}

bool AddGlobalsToSymbolTables(CompileState *state)
{
	AddGlobalsCtx ctx;
	ctx.ret = true;
	ctx.state = state;

	MapIterate(&state->root_namespace.children, AddGlobalsToNamespace, &ctx);

	return ctx.ret;
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

	if (!AddGlobalsToSymbolTables(state))
		ret = false;

	return ret;
}

