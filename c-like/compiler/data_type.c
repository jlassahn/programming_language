
#include "compiler/data_type.h"
#include "compiler/memory.h"
#include <stdlib.h>

DataType *DTypeMakeScalar(DTypeBase dts, uint32_t flags)
{
	DataType *dt = Alloc(sizeof(DataType));
	dt->base_type = dts;
	dt->flags = flags;
	return dt;
}

void DTypeFree(DataType *dt)
{
	// free subtypes etc
	Free(dt);
}

