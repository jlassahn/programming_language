
#include "compiler/symbol.h"
#include "compiler/memory.h"
#include <stdlib.h>

Symbol *SymbolCreate(const String *name)
{
	Symbol *sym = Alloc(sizeof(Symbol));
	sym->name = *name;

	return sym;
}

void SymbolDestroy(Symbol *sym)
{
	DValueClear(&sym->dvalue);

	if (sym->dtype != NULL)
		DTypeFree(sym->dtype);

	while (ListRemoveFirst(&sym->definitions) != NULL)
		; // don't free, owned by the compiler_file

	Free(sym);
}

DataValue *SymbolGetDValue(Symbol *sym)
{
	return &sym->dvalue;
}

DataType *SymbolGetDType(Symbol *sym)
{
	return sym->dtype;
}

