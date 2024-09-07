
#include "tests/unit/unit_test.h"
#include "compiler/data_value.h"
#include <stdbool.h>

void TestDataValue(void)
{
	DataValue dv;

	DataType *dt = DTypeMakeScalar(DTYPE_INT32, 0);

	DValueSetToDType(&dv, dt);
	CHECK(dv.value_type == VTYPE_DTYPE);
	CHECK(dv.value.dtype == dt);
	CHECK(dt->refcount == 2);
	DValueClear(&dv);
	CHECK(dv.value_type == VTYPE_INVALID);
	CHECK(dt->refcount == 1);

	DTypeFree(dt);
	// FIXME
	// data values can be
	//  - variable reference
	//  - constant number
	//  - constant large object (e.g. array or struct)
	//  - temporary value
	//  - data type
}

