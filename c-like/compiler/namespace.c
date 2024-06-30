
#include "compiler/namespace.h"
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/compiler_file.h"
#include "compiler/fileio.h"
#include <stdlib.h>

Namespace *NamespaceGetChild(Namespace *parent, String *name)
{
	Namespace *child = NULL;
	child = MapFind(&parent->children, name);
	if (!child)
	{
		child = Alloc(sizeof(Namespace));
		MapInsert(&parent->children, name, child);

		child->parent = parent;
		StringBuffer *sb = StringBufferFromString(&parent->path->string);
		int stem_offset = sb->string.length;
		sb = StringBufferAppendString(sb, name);
		sb = StringBufferAppendChars(sb, PATH_SEPARATOR_STRING);
		child->path = sb;
		child->stem.length = name->length;
		child->stem.data = sb->string.data + stem_offset;
	}

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
		CompilerFile *cf = ListRemoveFirst(&root->public_files);
		if (cf == NULL)
			break;

		// don't free the file if it's owned by the input list
		if (!(cf->flags & FILE_FROM_INPUT))
			CompilerFileFree(cf);
	}

	while (true)
	{
		CompilerFile *cf = ListRemoveFirst(&root->private_files);
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

	// FIXME free these...
	// Map public_symbols;  // FIXME Map of ????
	// Map private_symbols;  // FIXME Map of ????
}

void NamespacePrinter(const String *key, void *value, void *ctx)
{
	int *depth = ctx;
	Namespace *ns = value;

	for (int i=0; i<*depth; i++)
		printf("  ");
	printf("%.*s: %s\n", key->length, key->data, ns->path->buffer);

	*depth = *depth + 1;
	MapIterate(&ns->children, NamespacePrinter, depth);
	*depth = *depth - 1;
}

