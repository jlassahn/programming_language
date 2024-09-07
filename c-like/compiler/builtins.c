
#include "compiler/builtins.h"
#include "compiler/symbol.h"
#include <string.h>

static void AddDType(Map *map, const char *name, DTypeBase dts)
{
	String str;

	str.data = name;
	str.length = strlen(str.data);
	Symbol *sym = SymbolCreate(&str);

	DataType *dt = DTypeMakeScalar(dts, 0);
	DValueSetToDType(&sym->dvalue, dt);
	DTypeFree(dt);

	MapInsert(map, &str, sym);
}

bool InitBuiltins(Map *map, int bus_bits)
{

	// AddDType doesn't copy name, it must be
	// a long-lived char array.
	AddDType(map, "void", DTYPE_VOID);
	AddDType(map, "int", DTYPE_INT);
	AddDType(map, "float", DTYPE_FLOAT);
	AddDType(map, "int8", DTYPE_INT8);
	AddDType(map, "int16", DTYPE_INT16);
	AddDType(map, "int32", DTYPE_INT32);
	AddDType(map, "int64", DTYPE_INT64);
	AddDType(map, "uint8", DTYPE_UINT8);
	AddDType(map, "uint16", DTYPE_UINT16);
	AddDType(map, "uint32", DTYPE_UINT32);
	AddDType(map, "uint64", DTYPE_UINT64);
	AddDType(map, "float32", DTYPE_FLOAT32);
	AddDType(map, "float64", DTYPE_FLOAT64);

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


