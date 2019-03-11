
package symbols

var NAMESPACE_TYPE = &Tag{"NAMESPACE"}
var FUNCTIONCHOICE_TYPE = &Tag{"FUNCTIONCHOICE"}
var FUNCTION_TYPE = &Tag{"FUNCTION"}
var LABEL_TYPE = &Tag{"LABEL"}
var METATYPE_TYPE = &Tag{"METATYPE"}
var TYPEVAR_TYPE = &Tag{"TYPEVAR"}

var MREF_TYPE = &Tag{"MREF"}
var REF_TYPE = &Tag{"REF"}
var MSTRUCT_TYPE = &Tag{"MSTRUCT"}
var STRUCT_TYPE = &Tag{"STRUCT"}
var MARRAY_TYPE = &Tag{"MARRAY"}

var VOID_TYPE = &Tag{"VOID"}
var CTYPE_TYPE = &Tag{"CTYPE"}

var BOOL_TYPE = &Tag{"BOOL"}

var INTEGER_TYPE = &Tag{"INTEGER"}

var INT8_TYPE = &Tag{"INT8"}
var INT16_TYPE = &Tag{"INT16"}
var INT32_TYPE = &Tag{"INT32"}
var INT64_TYPE = &Tag{"INT64"}

var UINT8_TYPE = &Tag{"UINT8"}
var UINT16_TYPE = &Tag{"UINT16"}
var UINT32_TYPE = &Tag{"UINT32"}
var UINT64_TYPE = &Tag{"UINT64"}

var REAL32_TYPE = &Tag{"REAL32"}
var REAL64_TYPE = &Tag{"REAL64"}


var NamespaceType = &simpleDT{NAMESPACE_TYPE, nil}
var FunctionChoiceType = &simpleDT{FUNCTIONCHOICE_TYPE, nil}
var LabelType = &simpleDT{LABEL_TYPE, nil}
var MetaTypeType = &simpleDT{METATYPE_TYPE, nil}

var VoidType = &simpleDT{VOID_TYPE, nil}
var CTypeType = &simpleDT{CTYPE_TYPE, nil}

var BoolType = &simpleDT{BOOL_TYPE, nil}

var IntegerType = &simpleDT{INTEGER_TYPE, nil}
var Int8Type = &simpleDT{INT8_TYPE, nil}
var Int16Type = &simpleDT{INT16_TYPE, nil}
var Int32Type = &simpleDT{INT32_TYPE, nil}
var Int64Type = &simpleDT{INT64_TYPE, nil}
var UInt8Type = &simpleDT{UINT8_TYPE, nil}
var UInt16Type = &simpleDT{UINT16_TYPE, nil}
var UInt32Type = &simpleDT{UINT32_TYPE, nil}
var UInt64Type = &simpleDT{UINT64_TYPE, nil}

var Real32Type = &simpleDT{REAL32_TYPE, nil}
var Real64Type = &simpleDT{REAL64_TYPE, nil}

var TrueValue = &boolDV{ true }
var FalseValue = &boolDV{ false }

func InitializeTypes() {

	addTypeConvert("ToInt32", Int64Type, Int32Type)
	addTypeConvert("ToInt16", Int64Type, Int16Type)
	addTypeConvert("ToInt8", Int64Type, Int8Type)
	addTypeConvert("ToInt16", Int32Type, Int16Type)
	addTypeConvert("ToInt8", Int32Type, Int8Type)
	addTypeConvert("ToInt8", Int16Type, Int8Type)

}

func addTypeConvert(name string, from *simpleDT, to DataType) {

	if from.members == nil {
		from.members = map[string]Symbol { }
	}

	choices := &functionChoiceSymbol {
		name: name,
		choices: nil,
		modulePath: nil,
	}

	fnType := NewFunction(to)
	fnType.AddParam("a", from, false)
	choices.Add(
		&baseSymbol {
			name: name,
			dtype: fnType,
			initialValue: &intrinsicDV{fnType, "convert"},
			isConst: true,
			modulePath: nil,
			genVal: nil,
		})

	from.members[name] = choices
}

