
#include "compiler/compiler_file.h"
#include "compiler/memory.h"
#include "compiler/types.h"
#include "compiler/parser_file.h"
#include "compiler/parser.h"
#include "compiler/fileio.h"
#include "compiler/tokenizer.h"
#include <string.h>

CompilerFile *CompilerFileCreate(StringBuffer *path)
{
	CompilerFile *cf = Alloc(sizeof(CompilerFile));
	cf->path = path;
	return cf;
}

void CompilerFileFree(CompilerFile *cf)
{
	if (cf->path)
		StringBufferFree(cf->path);
	if (cf->parser_file)
		FileFree(cf->parser_file);
	if (cf->root)
		FreeNode(cf->root);

	Free(cf);
}

bool CompilerFilePickNamespace(CompilerFile *cf, Namespace *root)
{
	// For source files from the command line,
	// default behavior is to use the stem of the filename as the namespace.

	// FIXME official API for string search?
	const char *name = cf->path->buffer;
	int length = strlen(name);
	int stem_start = 0;
	for (int i=0; name[i]!=0; i++)
	{
		if (name[i] == PATH_SEPARATOR)
			stem_start = i+1;
	}
	int first_dot = length;
	for (int i=stem_start; name[i]!=0; i++)
	{
		if (name[i] == '.')
		{
			first_dot = i;
			break;
		}
	}

	if (first_dot - stem_start <= 0)
		return false;

	String stem;
	stem.data = name + stem_start;
	stem.length = first_dot - stem_start;

	if (!IsValidNamespaceName(&stem))
		return false;

	Namespace *namespace = NamespaceGetChild(root, &stem);
	ListInsertLast(&namespace->private_files, cf);
	namespace->flags |= NAMESPACE_HAS_INFILE;

	cf->namespace = namespace;

	return true;
}

