
package symbols

//FIXME Real or Float????

//FIXME do we need to expose bare tags for base types, or
//      just DataType values?

var VOID_TYPE = &Tag{"VOID"}
var NAMESPACE_TYPE = &Tag{"NAMESPACE"}
var FUNCTIONCHOICE_TYPE = &Tag{"FUNCTIONCHOICE"}
var FUNCTION_TYPE = &Tag{"FUNCTION"}
var INTRINSIC_TYPE = &Tag{"INTRINSIC"}
var CODE_TYPE = &Tag{"CODE"}
var LABEL_TYPE = &Tag{"LABEL"}

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


var VoidType = &simpleDT{VOID_TYPE}
var NamespaceType = &simpleDT{NAMESPACE_TYPE}
var FunctionChoiceType = &simpleDT{FUNCTIONCHOICE_TYPE}
var IntrinsicType = &simpleDT{INTRINSIC_TYPE}
var CodeType = &simpleDT{CODE_TYPE}
var LabelType = &simpleDT{LABEL_TYPE}

var CTypeType = &simpleDT{CTYPE_TYPE}

var BoolType = &simpleDT{BOOL_TYPE}

var IntegerType = &simpleDT{INTEGER_TYPE}
var Int8Type = &simpleDT{INT8_TYPE}
var Int16Type = &simpleDT{INT16_TYPE}
var Int32Type = &simpleDT{INT32_TYPE}
var Int64Type = &simpleDT{INT64_TYPE}
var UInt8Type = &simpleDT{UINT8_TYPE}
var UInt16Type = &simpleDT{UINT16_TYPE}
var UInt32Type = &simpleDT{UINT32_TYPE}
var UInt64Type = &simpleDT{UINT64_TYPE}

var Real32Type = &simpleDT{REAL32_TYPE}
var Real64Type = &simpleDT{REAL64_TYPE}

var TrueValue = &boolDV{ true }
var FalseValue = &boolDV{ false }

