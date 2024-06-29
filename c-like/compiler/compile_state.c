
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
}

