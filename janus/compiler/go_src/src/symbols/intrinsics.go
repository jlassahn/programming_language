
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
	"add_Int64": intrinsicAddInt64,
	"add_Int8": intrinsicAddInt,
	"div_Real64": intrinsicDivReal64,
}

var baseTypeConvert = map[tagPair] func(DataValue, DataType)DataValue  {
	{INT8_TYPE, INT16_TYPE}: convSignedSigned,
}


var intMask = map[*Tag] uint64 {
	INT8_TYPE  : 0x00000000000000FF,
	INT16_TYPE : 0x000000000000FFFF,
	INT32_TYPE : 0x00000000FFFFFFFF,
	INT64_TYPE : 0xFFFFFFFFFFFFFFFF,
	UINT8_TYPE : 0x00000000000000FF,
	UINT16_TYPE: 0x000000000000FFFF,
	UINT32_TYPE: 0x00000000FFFFFFFF,
	UINT64_TYPE: 0xFFFFFFFFFFFFFFFF,
}

type IntrinsicDataValue interface {
	ValueAsString() string
}

type intrinsicDV struct {
	name string
}

func (self *intrinsicDV) Type() DataType {
	return IntrinsicType
}

func (self *intrinsicDV) ValueAsString() string {
	return self.name
}

func (self *intrinsicDV) String() string {
	return DataValueString(self)
}

func intrinsicAddInt(op string, args []DataValue) DataValue {

	mask := intMask[args[0].Type().Base()]

	a := args[0].(*signedDV).AsSigned64()
	b := args[1].(*signedDV).AsSigned64()

	x := uint64(a+b) & mask
	if (x & (^mask >> 1)) != 0 {
		x = x | ^mask
	}

	return &signedDV{args[0].Type(), int64(x)}
}

func intrinsicAddInt64(op string, args []DataValue) DataValue {

	a := args[0].(*signedDV).AsSigned64()
	b := args[1].(*signedDV).AsSigned64()

	return &signedDV{Int64Type, a+b}
}

//FIXME should be real args, with conversions
func intrinsicDivReal64(op string, args []DataValue) DataValue {

	var a float64
	switch args[0].(type) {
		case *signedDV:
			a = float64(args[0].(*signedDV).AsSigned64())
		case *realDV:
			a = args[0].(*realDV).AsReal64()
	}

	var b float64
	switch args[1].(type) {
		case *signedDV:
			b = float64(args[1].(*signedDV).AsSigned64())
		case *realDV:
			b = args[1].(*realDV).AsReal64()
	}

	return &realDV{Real64Type, a/b}
}

func convSignedSigned(from DataValue, to DataType) DataValue {
	//FIXME mask
	return &signedDV{
		dtype: to,
		value: from.(SignedDataValue).AsSigned64(),
	}
}

