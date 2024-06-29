
#ifndef INCLUDED_COMPILE_STATE_H
#define INCLUDED_COMPILE_STATE_H

#include "compiler/types.h"

typedef struct CompileState CompileState;
struct CompileState
{
	List basedirs; // of locked StringBuffer*
};

void CompileStateFree(CompileState *state);

#endif

