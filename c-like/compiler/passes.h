
#ifndef INCLUDED_PASSES_H
#define INCLUDED_PASSES_H

#include "compiler/compile_state.h"
#include "compiler/commandargs.h"
#include <stdbool.h>

bool PassConfigure(CompileState *state, const CompilerArgs *args,
		const char *env);
bool PassSearchAndParse(CompileState *state);
bool PassResolveGlobals(CompileState *state);

#endif

