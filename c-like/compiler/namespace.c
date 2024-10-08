
#include "compiler/namespace.h"
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/compiler_file.h"
#include "compiler/fileio.h"
#include <stdlib.h>

Namespace *NamespaceMakeChild(Namespace *parent, String *name)
{
	Namespace *child = NULL;
	child = MapFind(&parent->children, name);
	if (!child)
	{
		child = Alloc(sizeof(Namespace));

		child->parent = parent;
		StringBuffer *sb = StringBufferFromString(&parent->path->string);
		int stem_offset = sb->string.length;
		sb = StringBufferAppendString(sb, name);
		sb = StringBufferAppendChars(sb, PATH_SEPARATOR_STRING);
		child->path = sb;
		child->stem.length = name->length;
		child->stem.data = sb->string.data + stem_offset;

		MapInsert(&parent->children, &child->stem, child);
	}

	return child;
}

Namespace *NamespaceGetChild(Namespace *parent, String *name)
{
	Namespace *child = NULL;
	child = MapFind(&parent->children, name);
	return child;
}

void NamespaceFree(Namespace *root)
{
	root->flags = 0;
	root->parent = NULL;
	if (root->path)
	{
		StringBufferFree(root->path);
		root->path = NULL;
	}
	root->stem.length = 0;
	root->stem.data = NULL;

	while (true)
	{
		CompilerFile *cf = ListRemoveFirst(&root->public_syms.files);
		if (cf == NULL)
			break;

		// don't free the file if it's owned by the input list
		if (!(cf->flags & FILE_FROM_INPUT))
			CompilerFileFree(cf);
	}

	while (true)
	{
		CompilerFile *cf = ListRemoveFirst(&root->private_syms.files);
		if (cf == NULL)
			break;

		// don't free the file if it's owned by the input list
		if (!(cf->flags & FILE_FROM_INPUT))
			CompilerFileFree(cf);
	}

	while (true)
	{
		Namespace *child = MapRemoveFirst(&root->children);
		if (child == NULL)
			break;
		NamespaceFree(child);
		Free(child);
	}

	while (true)
	{
		Symbol *sym = MapRemoveFirst(&root->public_syms.exports);
		if (sym == NULL)
			break;
		SymbolDestroy(sym);
	}

	while (true)
	{
		Symbol *sym = MapRemoveFirst(&root->private_syms.exports);
		if (sym == NULL)
			break;
		SymbolDestroy(sym);
	}

	// FIXME untangle ownership of ImportLink
	while (true)
	{
		ImportLink *import = ListRemoveFirst(&root->public_syms.imports);
		if (import == NULL)
			break;
		// don't destroy import here, all of these are also in private_syms.imports
	}

	while (true)
	{
		ImportLink *import = ListRemoveFirst(&root->private_syms.imports);
		if (import == NULL)
			break;
		// currently destroyed by CompilerFile
	}

	SymbolTableDestroy(&root->public_syms.symbol_table);
	SymbolTableDestroy(&root->private_syms.symbol_table);
}

Symbol *NamespaceFindSymbol(Namespace *ns, String *name)
{
	return MapFind(&ns->public_syms.exports, name);
}

Symbol *NamespaceFindPrivateSymbol(Namespace *ns, String *name)
{
	return MapFind(&ns->private_syms.exports, name);
}

void SymbolPrinter(const String *key, void *value, void *ctx)
{
	int *depth = ctx;
	// Symbol *sym = value;

	for (int i=0; i<*depth; i++)
		printf("  ");
	printf("-> %.*s\n", key->length, key->data);
}

void NamespacePrinter(const String *key, void *value, void *ctx)
{
	int *depth = ctx;
	Namespace *ns = value;

	for (int i=0; i<*depth; i++)
		printf("  ");
	printf("%.*s: %s\n", key->length, key->data, ns->path->buffer);

	*depth = *depth + 1;
	MapIterate(&ns->private_syms.exports, SymbolPrinter, depth);
	MapIterate(&ns->children, NamespacePrinter, depth);
	*depth = *depth - 1;
}

