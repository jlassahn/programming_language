
#ifndef INCLUDED_COMPILE_STATE_H
#define INCLUDED_COMPILE_STATE_H

#include "compiler/types.h"
#include "compiler/namespace.h"

typedef struct CompileState CompileState;
struct CompileState
{
	List basedirs; // List of locked StringBuffer*
	List input_files; // List of CompilerFile*
	List input_modules; // List of Namespace*

	Namespace root_namespace;
};

void CompileStateInit(CompileState *state);
void CompileStateFree(CompileState *state);
void CompileStatePrint(const CompileState *state);

#endif

