
#include "tests/unit/unit_test.h"
#include "compiler/data_type.h"
#include <stdbool.h>

void TestDataType(void)
{
	DataType *dt = DTypeMakeScalar(DTYPE_UINT32, TYPEFLAG_READONLY);
	CHECK(dt->base_type == DTYPE_UINT32);
	CHECK(dt->flags == TYPEFLAG_READONLY);
	CHECK(dt->refcount == 1);

	DataType *dt2 = DTypeCopy(dt);
	CHECK(dt == dt2);
	CHECK(dt->refcount == 2);

	DTypeFree(dt);
	CHECK(dt2->refcount == 1);
	DTypeFree(dt2);
}

