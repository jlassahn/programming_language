
#include "compiler/compile_state.h"
#include "compiler/types.h"
#include "compiler/compiler_file.h"
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
		StringBuffer *sb = ListRemoveFirst(&state->input_modules);
		if (sb == NULL)
			break;
		StringBufferFree(sb);
	}

	NamespaceFree(&state->root_namespace);
}

