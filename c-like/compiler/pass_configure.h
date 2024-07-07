
#ifndef INCLUDED_PASS_CONFIGURE_H
#define INCLUDED_PASS_CONFIGURE_H

#include "compiler/compile_state.h"
#include "compiler/commandargs.h"
#include <stdbool.h>

bool PassConfigure(CompileState *state, const CompilerArgs *args,
		const char *env);

#endif

