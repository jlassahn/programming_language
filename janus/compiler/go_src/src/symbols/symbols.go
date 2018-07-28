
package symbols

type Tag struct { string }

var VOID_TYPE = &Tag{"VOID"}
var BOOL_TYPE = &Tag{"BOOL"}
var INT_TYPE = &Tag{"INT"}


type DataType struct {
	BaseType *Tag
	SubTypes []*DataType
}

type DataValue struct {
	Type *DataType
	Int64 int64
	UInt64 uint64
	Float64 float64
}

type Symbol struct {
	Name string
	Type *DataType
	InitialValue *DataValue
}

type SymbolTable struct {
	Symbols map[string]*Symbol
	Parent *SymbolTable
}

func ResolveGlobals(file_set *FileSet) {
	//FIXME implement
}

var boolType = &DataType{BaseType: BOOL_TYPE}

var trueValue = &DataValue{ Type: boolType }
var falseValue = &DataValue{ Type: boolType }


//FIXME implement
var PredefinedSymbols = &SymbolTable {
	Symbols: map[string]*Symbol {
		"True": &Symbol {"True", boolType, trueValue } ,
		"False": &Symbol {"False", boolType, falseValue } },
	Parent: nil }

