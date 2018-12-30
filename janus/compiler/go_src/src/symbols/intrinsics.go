
package symbols

import (
	"output"
)

func EvaluateIntrinsic(opName string, args []DataValue) DataValue {

	fn := intrinsics[opName]
	if fn == nil {
		output.FatalError("Unimplemented intrinsic: %v", opName)
		return nil
	}

	return fn(opName, args)
}

func ConvertBasic(from DataValue, to DataType) DataValue {

	fn := baseTypeConvert[tagPair{from.Type().Base(), to.Base()}]
	if fn == nil {
		return nil
	}

	return fn(from, to)
}

var intrinsics = map[string] func(name string, args []DataValue) DataValue {
	"add_Int8": intrinsicAddInt,
	"add_Int16": intrinsicAddInt,
	"add_Int32": intrinsicAddInt,
	"add_Int64": intrinsicAddInt,
	"add_Real32": intrinsicAddReal,
	"add_Real64": intrinsicAddReal,

	"sub_Int8": intrinsicSubInt,
	"sub_Int16": intrinsicSubInt,
	"sub_Int32": intrinsicSubInt,
	"sub_Int64": intrinsicSubInt,
	"sub_Real32": intrinsicSubReal,
	"sub_Real64": intrinsicSubReal,

	"negate_Int8": intrinsicNegateInt,
	"negate_Int16": intrinsicNegateInt,
	"negate_Int32": intrinsicNegateInt,
	"negate_Int64": intrinsicNegateInt,
	"negate_Real32": intrinsicNegateReal,
	"negate_Real64": intrinsicNegateReal,

	"div_IntReal": intrinsicDivIntReal,
	"div_Real64": intrinsicDivReal,
	"div_Real32": intrinsicDivReal,

	//FIXME
	//"cmp_le_Int64": intrinsicCmpLEInt,
}

var baseTypeConvert = map[tagPair] func(DataValue, DataType)DataValue  {
	{INT8_TYPE, INT16_TYPE}: convSignedSigned,
	{INT8_TYPE, INT32_TYPE}: convSignedSigned,
	{INT8_TYPE, INT64_TYPE}: convSignedSigned,
	{INT16_TYPE, INT32_TYPE}: convSignedSigned,
	{INT16_TYPE, INT64_TYPE}: convSignedSigned,
	{INT32_TYPE, INT64_TYPE}: convSignedSigned,

	{INT8_TYPE, REAL32_TYPE}: convSignedReal,
	{INT16_TYPE, REAL32_TYPE}: convSignedReal,
	{INT32_TYPE, REAL32_TYPE}: convSignedReal,
	{INT64_TYPE, REAL32_TYPE}: convSignedReal,
	{INT8_TYPE, REAL64_TYPE}: convSignedReal,
	{INT16_TYPE, REAL64_TYPE}: convSignedReal,
	{INT32_TYPE, REAL64_TYPE}: convSignedReal,
	{INT64_TYPE, REAL64_TYPE}: convSignedReal,

}


var intMask = map[*Tag] uint64 {
	INT8_TYPE  : 0x00000000000000FF,
	INT16_TYPE : 0x000000000000FFFF,
	INT32_TYPE : 0x00000000FFFFFFFF,
	INT64_TYPE : 0xFFFFFFFFFFFFFFFF,
}

var uintMask = map[*Tag] uint64 {
	UINT8_TYPE : 0x00000000000000FF,
	UINT16_TYPE: 0x000000000000FFFF,
	UINT32_TYPE: 0x00000000FFFFFFFF,
	UINT64_TYPE: 0xFFFFFFFFFFFFFFFF,
}

type IntrinsicDataValue interface {
	ValueAsString() string
	IsIntrinsic() bool //FIXME remove!
}

type intrinsicDV struct {
	dtype DataType
	name string
}

func (self *intrinsicDV) Tag() *ValueTag {
	return INTRINSIC_VALUE
}

func (self *intrinsicDV) Type() DataType {
	return self.dtype
}

func (self *intrinsicDV) ValueAsString() string {
	return self.name
}

func (self *intrinsicDV) IsIntrinsic() bool {
	return true
}

func (self *intrinsicDV) String() string {
	return DataValueString(self)
}

func maskSigned(x int64, dtype *Tag) int64 {

	mask := intMask[dtype]

	m := uint64(x) & mask
	if (m & (^mask >> 1)) != 0 {
		m = m | ^mask
	}
	ret := int64(m)

	//FIXME better context for message!
	if x != ret {
		output.Warning("truncating constant value %v to %v", x, ret)
	}

	return ret
}

func intrinsicAddInt(op string, args []DataValue) DataValue {

	a := args[0].(*signedDV).AsSigned64()
	b := args[1].(*signedDV).AsSigned64()

	x := maskSigned(a+b, args[0].Type().Base())

	return &signedDV{args[0].Type(), x}
}

func intrinsicAddReal(op string, args []DataValue) DataValue {

	a := args[0].(*realDV).AsReal64()
	b := args[1].(*realDV).AsReal64()

	x := a+b

	return &realDV{args[0].Type(), x}
}

func intrinsicSubInt(op string, args []DataValue) DataValue {

	a := args[0].(*signedDV).AsSigned64()
	b := args[1].(*signedDV).AsSigned64()

	x := maskSigned(a-b, args[0].Type().Base())

	return &signedDV{args[0].Type(), x}
}

func intrinsicSubReal(op string, args []DataValue) DataValue {

	a := args[0].(*realDV).AsReal64()
	b := args[1].(*realDV).AsReal64()

	x := a-b

	return &realDV{args[0].Type(), x}
}

func intrinsicNegateInt(op string, args []DataValue) DataValue {

	a := args[0].(*signedDV).AsSigned64()

	x := maskSigned(-a, args[0].Type().Base())

	return &signedDV{args[0].Type(), x}
}

func intrinsicNegateReal(op string, args []DataValue) DataValue {

	a := args[0].(*realDV).AsReal64()

	x := -a

	return &realDV{args[0].Type(), x}
}

func intrinsicDivIntReal(op string, args []DataValue) DataValue {

	a := float64(args[0].(*signedDV).AsSigned64())
	b := args[1].(*realDV).AsReal64()
	return &realDV{Real64Type, a/b}
}

func intrinsicDivReal(op string, args []DataValue) DataValue {

	a := args[0].(*realDV).AsReal64()
	b := args[1].(*realDV).AsReal64()
	return &realDV{args[0].Type(), a/b}
}

func convSignedSigned(from DataValue, to DataType) DataValue {

	x := from.(SignedDataValue).AsSigned64()
	x = maskSigned(x, from.Type().Base())
	x = maskSigned(x, to.Base())

	return &signedDV{
		dtype: to,
		value: x,
	}
}

func convSignedReal(from DataValue, to DataType) DataValue {

	x := from.(SignedDataValue).AsSigned64()
	x = maskSigned(x, from.Type().Base())

	return &realDV{
		dtype: to,
		value: float64(x),
	}
}

func MaskConstant(x DataValue) DataValue {

	mask := intMask[x.Type().Base()]
	if mask != 0 {

		inVal := x.(SignedDataValue).AsSigned64()
		outVal := maskSigned(inVal, x.Type().Base())
		if inVal == outVal {
			return x
		}

		return &signedDV{
			dtype: x.Type(),
			value: outVal,
		}
	}

	//FIXME implement unsigned

	return x
}

