

#include "tests/unit/unit_test.h"
#include "tests/unit/utils.h"
#include "compiler/types.h"
#include "compiler/builtins.h"
#include <string.h>

void TestBuiltins(void)
{
	Map map;
	memset(&map, 0, sizeof(Map));

	CHECK(InitBuiltins(&map, 32));
	CHECK(MapFind(&map, TempString("void")) != NULL);
	CHECK(MapFind(&map, TempString("int")) != NULL);
	CHECK(MapFind(&map, TempString("float")) != NULL);
	CHECK(MapFind(&map, TempString("int8")) != NULL);
	CHECK(MapFind(&map, TempString("int16")) != NULL);
	CHECK(MapFind(&map, TempString("int32")) != NULL);
	CHECK(MapFind(&map, TempString("int64")) != NULL);
	CHECK(MapFind(&map, TempString("uint8")) != NULL);
	CHECK(MapFind(&map, TempString("uint16")) != NULL);
	CHECK(MapFind(&map, TempString("uint32")) != NULL);
	CHECK(MapFind(&map, TempString("uint64")) != NULL);
	CHECK(MapFind(&map, TempString("float32")) != NULL);
	CHECK(MapFind(&map, TempString("float64")) != NULL);

	FreeBuiltins(&map);

}

