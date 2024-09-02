

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

	FreeBuiltins(&map);

}

