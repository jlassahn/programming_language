
#ifndef INCLUDED_DATA_VALUE_H
#define INCLUDED_DATA_VALUE_H

#include "compiler/data_type.h"

typedef enum
{
	VTYPE_INVALID,
	VTYPE_DTYPE
}
ValueType;

typedef union ValueUnion ValueUnion;
union ValueUnion
{
	DataType *dtype;
};

typedef struct DataValue DataValue;
struct DataValue
{
	ValueType value_type;
	ValueUnion value;
};

void DValueSetToDType(DataValue *dv, DataType *dt);
void DValueClear(DataValue *dv);

#endif

