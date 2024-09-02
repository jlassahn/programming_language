
#ifndef INCLUDED_BUILTINS_H
#define INCLUDED_BUILTINS_H

#include "compiler/types.h"
#include <stdbool.h>

bool InitBuiltins(Map *map, int bus_bits);
void FreeBuiltins(Map *map);

#endif

