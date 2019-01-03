
package symbols

import (
	"output"
)

/* FIXME implicit type conversion rules

These are permitted:
 Int8 -> Int16 -> Int32 -> Int64 -> Integer -> Real32 -> Real64
 UInt8 -> UInt16 -> UInt32 -> UInt64
 MArray -> MSlice
 MArray -> Array
(Can't do MSlice -> Array or Array -> others automatically)
Any combination of Ref and MRef can be removed from outermost
	(Anything passed by Ref must match exactly in base type)
One level of Ref or MRef can be added (to outermost nonref)
	(Must be a variable of exactly matching type)
(Can't both remove and add, i.e. can't change between Ref and MRef)
 Structs convert to anything they extend.
 Structs have to have the same named type, just having matching elements
    isn't sufficient.
 Structs and interfaces can convert to anything they implement
 Functions must have the same param and return types
 Methods can become functions with this as the first parameter

Methods collide when parameters have these properties:
	They match
	an autoconvert param can implicitly convert to a non-autoconvert
	Both autoconvert and
		either one is an interface of any kind
			(can build a new type that extends a and implements b)
		either one can convert to the other
			(non-interface types form a tree, so paths between them are
			always with one as the deepest and one as the shallowest
			point on the path)
*/

func CanConvertArgs(args []DataType, params []FunctionParameter) bool {

	if len(args) != len(params) {
		return false
	}

	for i, param := range params {
		arg := args[i]

		if param.AutoConvert {
			if !CanConvert(arg, param.DType) {
				return false
			}
		} else {
			if !TypeMatches(arg, param.DType) {
				return false
			}
		}
	}

	return true
}


type tagPair struct {
	from *Tag
	to *Tag
}

var baseTypeMatch = map[tagPair] bool {
	{INT8_TYPE, INT16_TYPE} : true,
	{INT8_TYPE, INT32_TYPE} : true,
	{INT8_TYPE, INT64_TYPE} : true,
	{INT8_TYPE, REAL32_TYPE} : true,
	{INT8_TYPE, REAL64_TYPE} : true,

	{INT16_TYPE, INT32_TYPE} : true,
	{INT16_TYPE, INT64_TYPE} : true,
	{INT16_TYPE, REAL32_TYPE} : true,
	{INT16_TYPE, REAL64_TYPE} : true,

	{INT32_TYPE, INT64_TYPE} : true,
	{INT32_TYPE, REAL32_TYPE} : true,
	{INT32_TYPE, REAL64_TYPE} : true,

	{INT64_TYPE, REAL32_TYPE} : true,
	{INT64_TYPE, REAL64_TYPE} : true,

	{REAL32_TYPE, REAL64_TYPE} : true,

	{UINT8_TYPE, UINT16_TYPE} : true,
	{UINT8_TYPE, UINT32_TYPE} : true,
	{UINT8_TYPE, UINT64_TYPE} : true,

	{UINT16_TYPE, UINT32_TYPE} : true,
	{UINT16_TYPE, UINT64_TYPE} : true,

	{UINT32_TYPE, UINT64_TYPE} : true,
}

func CanConvert(argType DataType, paramType DataType) bool {

	if TypeMatches(argType, paramType) {
		return true
	}

	if baseTypeMatch[tagPair{argType.Base(), paramType.Base()}] {
		return true
	}

	//FIXME check composite types, refs, etc
	if argType.Base() == MREF_TYPE {
		output.FatalError("FIXME implement MREF conversions")
	}
	if argType.Base() == REF_TYPE {
		output.FatalError("FIXME implement REF conversions")
	}
	if argType.Base() == MSTRUCT_TYPE {
		output.FatalError("FIXME implement MSTRUCT conversions")
	}
	if argType.Base() == STRUCT_TYPE {
		output.FatalError("FIXME implement STRUCT conversions")
	}
	if argType.Base() == MARRAY_TYPE {
		output.FatalError("FIXME implement MARRAY conversions")
	}

	return false
}


func TypeMatches(a DataType, b DataType) bool {
	if a == b {
		return true
	}

	if a == nil {
		return false
	}

	if b == nil {
		return false
	}

	if a.Base() != b.Base() {
		return false
	}

	subA := a.SubTypes()
	subB := b.SubTypes()
	if len(subA) != len(subB) {
		return false
	}

	for i:=0; i<len(subA); i++ {
		if subA[i].Number != subB[i].Number {
			return false
		}

		if !TypeMatches(subA[i].DType, subB[i].DType) {
			return false
		}
	}

	if a.Base() == FUNCTION_TYPE {
		fnA := a.(FunctionDataType)
		fnB := b.(FunctionDataType)

		if fnA.IsMethod() != fnB.IsMethod() {
			return false
		}

		if !TypeMatches(fnA.ReturnType(), fnB.ReturnType()) {
			return false
		}

		paramsA := fnA.Parameters()
		paramsB := fnB.Parameters()
		if len(paramsA) != len(paramsB) {
			return false
		}
		for i:=0; i<len(paramsA); i++ {
			if paramsA[i].AutoConvert != paramsB[i].AutoConvert {
				return false
			}

			if !TypeMatches(paramsA[i].DType, paramsB[i].DType) {
				return false
			}
		}
	}

	//FIXME handle struct types
	if a.Base() == MSTRUCT_TYPE {
		output.FatalError("FIXME implement MSTRUCT conversions")
	}
	if a.Base() == STRUCT_TYPE {
		output.FatalError("FIXME implement STRUCT conversions")
	}

	return true
}

func ConvertConstant(from DataValue, to DataType) DataValue {

	if TypeMatches(from.Type(), to) {
		return from
	}

	ret := ConvertBasic(from, to)
	if ret != nil {
		return ret
	}

	if from.Type() == FunctionChoiceType {
		fn := from.(FunctionChoiceValue).AsSymbol()
		for _,choice := range fn.Choices() {
			if TypeMatches(choice.Type(), to) {

				ret := choice.InitialValue()

				//special case for assigning the value of a function
				// definition to another constant value, supply a global
				// data reference instead of a copy of the code.
				if ret.Tag() == CODE_VALUE {
					ret = &globalDV {
						dtype: choice.Type(),
						symbol: choice,
						offset: 0,
					}
				}

				return ret
			}
		}
	}

	//FIXME handle composite types, etc
	if from.Type().Base() == MREF_TYPE {
		output.FatalError("FIXME implement MREF conversions")
	}
	if from.Type().Base() == REF_TYPE {
		output.FatalError("FIXME implement REF conversions")
	}
	if from.Type().Base() == MSTRUCT_TYPE {
		output.FatalError("FIXME implement MSTRUCT conversions")
	}
	if from.Type().Base() == STRUCT_TYPE {
		output.FatalError("FIXME implement STRUCT conversions")
	}
	if from.Type().Base() == MARRAY_TYPE {
		output.FatalError("FIXME implement MARRAY conversions")
	}

	//FIXME better context for errors
	output.Error("no const conversion from %v to %v", from, to)
	return nil
}

