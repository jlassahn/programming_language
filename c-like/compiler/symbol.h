
#ifndef INCLUDED_SYMBOL_H
#define INCLUDED_SYMBOL_H

#include "compiler/types.h"
#include "compiler/data_value.h"

typedef struct Symbol Symbol;
struct Symbol
{
	String name;
	Symbol *associated; // link between public and private views of the symbol
	DataValue dvalue;
};

Symbol *SymbolCreate(const String *name);
void SymbolDestroy(Symbol *sym);
DataValue *SymbolGetDValue(Symbol *sym);

#endif

