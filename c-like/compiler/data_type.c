
#include "compiler/data_type.h"
#include "compiler/memory.h"
#include <stdlib.h>

DataType *DTypeMakeScalar(DTypeBase dts, uint32_t flags)
{
	DataType *dt = Alloc(sizeof(DataType));
	dt->base_type = dts;
	dt->flags = flags;
	dt->refcount = 1;
	return dt;
}

DataType *DTypeCopy(DataType *dt)
{
	dt->refcount ++;
	return dt;
}

void DTypeFree(DataType *dt)
{
	dt->refcount --;
	if (dt->refcount > 0)
		return;

	// free subtypes etc
	Free(dt);
}

