
package symbols

//FIXME Real or Float????

//FIXME do we need to expose bare tags for base types, or
//      just DataType values?

var VOID_TYPE = &Tag{"VOID"}
var NAMESPACE_TYPE = &Tag{"NAMESPACE"}

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


var VoidType = &SimpleDataType{*VOID_TYPE}
var NamespaceType = &SimpleDataType{*NAMESPACE_TYPE}
var BoolType = &SimpleDataType{*BOOL_TYPE}

var IntegerType = &SimpleDataType{*INTEGER_TYPE}
var Int8Type = &SimpleDataType{*INT8_TYPE}
var Int16Type = &SimpleDataType{*INT16_TYPE}
var Int32Type = &SimpleDataType{*INT32_TYPE}
var Int64Type = &SimpleDataType{*INT64_TYPE}
var UInt8Type = &SimpleDataType{*UINT8_TYPE}
var UInt16Type = &SimpleDataType{*UINT16_TYPE}
var UInt32Type = &SimpleDataType{*UINT32_TYPE}
var UInt64Type = &SimpleDataType{*UINT64_TYPE}

var Real32Type = &SimpleDataType{*REAL32_TYPE}
var Real64Type = &SimpleDataType{*REAL64_TYPE}

var TrueValue = &boolDV{ true }
var FalseValue = &boolDV{ false }

