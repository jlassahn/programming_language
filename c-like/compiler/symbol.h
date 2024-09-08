
#ifndef INCLUDED_SYMBOL_H
#define INCLUDED_SYMBOL_H

#include "compiler/types.h"
#include "compiler/data_value.h"

typedef struct Namespace Namespace;

typedef enum
{
	SYM_PRIVATE = 0x0001,
}
SymbolFlags;

typedef struct Symbol Symbol;
struct Symbol
{
	uint32_t flags;
	String name;
	Symbol *associated; // link between public and private views of the symbol
	Namespace *exported_from;
	DataValue dvalue;
	DataType *dtype;
	List definitions; // List of ParserNode
};

Symbol *SymbolCreate(const String *name);
void SymbolDestroy(Symbol *sym);
DataValue *SymbolGetDValue(Symbol *sym);
DataType *SymbolGetDType(Symbol *sym);

#endif

