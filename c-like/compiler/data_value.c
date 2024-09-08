
#include "compiler/data_value.h"
#include "compiler/errors.h"
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

void DValueCopy(DataValue *dest, DataValue *src)
{
	switch (src->value_type)
	{
		case VTYPE_DTYPE:
			dest->value.dtype = DTypeCopy(src->value.dtype);
			break;
		default:
			Error(ERROR_INTERNAL, "UNIMPLEMENTED data value copy on type %d", src->value_type);
			dest->value_type = VTYPE_INVALID;
			return;
	}

	dest->value_type = src->value_type;
}

