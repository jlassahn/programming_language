
#include "tests/unit/unit_test.h"
#include "tests/unit/utils.h"
#include "compiler/symbol_table.h"
#include "compiler/symbol.h"
#include "compiler/memory.h"
#include <string.h>

Map *CreateMap(void)
{
	Symbol *sym;
	Map *map = Alloc(sizeof(Map));
	sym = SymbolCreate(TempString("map3"));
	MapInsert(map, &sym->name, sym);
	sym = SymbolCreate(TempString("map2"));
	MapInsert(map, &sym->name, sym);
	sym = SymbolCreate(TempString("map1"));
	MapInsert(map, &sym->name, sym);
	return map;
}

void DestroyMap(Map *map)
{
	while (true)
	{
		Symbol *sym =MapRemoveFirst(map);
		if (sym == NULL)
			break;
		SymbolDestroy(sym);
	}
	Free(map);
}

void TestSymbolTable(void)
{
	SymbolTable syms;
	SymbolTableInit(&syms);

	CHECK(SymbolTableFind(&syms, TempString("test_symbol")) == NULL);
	Symbol *sym = SymbolCreate(TempString("test_symbol"));

	CHECK(SymbolTableInsert(&syms, sym));
	CHECK(!SymbolTableInsert(&syms, sym)); // second time should fail
	CHECK(SymbolTableFind(&syms, TempString("test_symbol")) == sym);

	Map *map = CreateMap();
	CHECK(SymbolTableInsertMap(&syms, map));
	CHECK(SymbolTableFind(&syms, TempString("test_symbol")) == sym);
	CHECK(SymbolTableFind(&syms, TempString("map1")) != NULL);
	CHECK(SymbolTableFind(&syms, TempString("map2")) != NULL);
	CHECK(SymbolTableFind(&syms, TempString("map3")) != NULL);

	SymbolTableDestroy(&syms);
	DestroyMap(map);
	SymbolDestroy(sym);
}

