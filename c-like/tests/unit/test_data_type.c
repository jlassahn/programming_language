
#include "tests/unit/unit_test.h"
#include "compiler/data_type.h"
#include <stdbool.h>

void TestDataType(void)
{
	DataType *dt = DTypeMakeScalar(DTYPE_UINT32, TYPEFLAG_READONLY);
	CHECK(dt->base_type == DTYPE_UINT32);
	CHECK(dt->flags == TYPEFLAG_READONLY);
	DTypeFree(dt);
}

