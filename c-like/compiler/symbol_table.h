
#ifndef INCLUDED_SYMBOL_TABLE_H
#define INCLUDED_SYMBOL_TABLE_H

#include "compiler/symbol.h"
#include "compiler/types.h"

typedef struct SymbolTable SymbolTable;
struct SymbolTable
{
	Map root; // Map of Symbol
};

void SymbolTableInit(SymbolTable *syms);
void SymbolTableDestroy(SymbolTable *syms);
Symbol *SymbolTableFind(SymbolTable *syms, String *name);
bool SymbolTableInsert(SymbolTable *syms, Symbol *sym);
bool SymbolTableInsertMap(SymbolTable *syms, Map *map);

#endif

