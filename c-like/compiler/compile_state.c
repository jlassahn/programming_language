
#include "compiler/compile_state.h"
#include "compiler/types.h"
#include <stdlib.h>

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
		StringBuffer *sb = ListRemoveFirst(&state->input_files);
		if (sb == NULL)
			break;
		StringBufferFree(sb);
	}

	while (true)
	{
		StringBuffer *sb = ListRemoveFirst(&state->input_modules);
		if (sb == NULL)
			break;
		StringBufferFree(sb);
	}
}

