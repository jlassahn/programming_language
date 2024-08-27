
#ifndef INCLUDED_DATA_TYPE_H
#define INCLUDED_DATA_TYPE_H

#include <stdint.h>
#include <stdbool.h>

typedef enum
{
	DTYPE_VOID,
	DTYPE_INT,
	DTYPE_FLOAT,
	DTYPE_INT8,
	DTYPE_INT16,
	DTYPE_INT32,
	DTYPE_INT64,
	DTYPE_UINT8,
	DTYPE_UINT16,
	DTYPE_UINT32,
	DTYPE_UINT64,
	DTYPE_FLOAT32,
	DTYPE_FLOAT64
}
DTypeBase;

typedef enum
{
	TYPEFLAG_READONLY = 0x0001
}
DTypeFlags;

typedef struct DataType DataType;
struct DataType
{
	DTypeBase base_type;
	uint32_t flags;
	// DataType *subtype;
	// union dtype_params;
};


DataType *DTypeMakeScalar(DTypeBase dts, uint32_t flags);
void DTypeFree(DataType *dt);

#endif

