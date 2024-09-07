
#include "compiler/compile_state.h"
#include "compiler/types.h"
#include "compiler/compiler_file.h"
#include "compiler/builtins.h"
#include <stdlib.h>
#include <string.h>

void CompileStateInit(CompileState *state)
{
	memset(state, 0, sizeof(CompileState));

	state->root_namespace.parent = NULL;
	state->root_namespace.path = StringBufferFromChars("");
	state->root_namespace.stem.length = 0;
	state->root_namespace.stem.data = "";
}

void CompileStateFree(CompileState *state)
{
	while (true)
	{
		StringBuffer *sb = ListRemoveFirst(&state->basedirs);
		if (sb == NULL)
			break;
		StringBufferFree(sb);
	}

	while (true)
	{
		CompilerFile *cf = ListRemoveFirst(&state->input_files);
		if (cf == NULL)
			break;
		CompilerFileFree(cf);
	}

	while (true)
	{
		Namespace *ns = ListRemoveFirst(&state->input_modules);
		if (ns == NULL)
			break;
		// Namespaces are owned by the root namespace below
	}

	NamespaceFree(&state->root_namespace);
	FreeBuiltins(&state->builtins);
}

void CompileStatePrint(const CompileState *state)
{
	for (ListEntry *entry=state->basedirs.first;
			entry!=NULL; entry=entry->next)
	{
		StringBuffer *sb = entry->item;
		printf("search directory: %s\n", sb->string.data);
	}

	for (ListEntry *entry=state->input_files.first;
			entry!=NULL; entry=entry->next)
	{
		CompilerFile *cf = entry->item;
		printf("input file: %s\n", cf->path->string.data);
	}

	for (ListEntry *entry=state->input_modules.first;
			entry!=NULL; entry=entry->next)
	{
		Namespace *ns = entry->item;
		printf("input module: %s\n", ns->path->buffer);
	}
}

