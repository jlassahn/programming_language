
#include "compiler/symbol.h"
#include "compiler/memory.h"

Symbol *SymbolCreate(const String *name)
{
	Symbol *sym = Alloc(sizeof(Symbol));
	sym->name = *name;

	return sym;
}

void SymbolDestroy(Symbol *sym)
{
	Free(sym);
}

