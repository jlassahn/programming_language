
#include "compiler/builtins.h"
#include "compiler/symbol.h"
#include <string.h>

bool InitBuiltins(Map *map, int bus_bits)
{
	String name;

	name.data = "void";
	name.length = strlen(name.data);
	Symbol *sym = SymbolCreate(&name);
	MapInsert(map, &name, sym);

	return true;
}

void FreeBuiltins(Map *map)
{

	while (true)
	{
		Symbol *sym = MapRemoveFirst(map);
		if (sym == NULL)
			break;
		SymbolDestroy(sym);
	}
}


