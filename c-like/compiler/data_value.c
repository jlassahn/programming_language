
#include "compiler/data_value.h"
#include <stdlib.h>

void DValueSetToDType(DataValue *dv, DataType *dt)
{
	dv->value_type = VTYPE_DTYPE;
	dv->value.dtype = DTypeCopy(dt);
}

void DValueClear(DataValue *dv)
{
	switch (dv->value_type)
	{
		case VTYPE_DTYPE:
			DTypeFree(dv->value.dtype);
			dv->value.dtype = NULL;
			break;
		default:
			break;
	}

	dv->value_type = VTYPE_INVALID;
}

