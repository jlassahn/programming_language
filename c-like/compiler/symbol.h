
#ifndef INCLUDED_SYMBOL_H
#define INCLUDED_SYMBOL_H

#include "compiler/types.h"

typedef struct Symbol Symbol;
struct Symbol
{
	String name;
};

Symbol *SymbolCreate(const String *name);
void SymbolDestroy(Symbol *sym);

#endif

