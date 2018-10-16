
package symbols

import (
	"fmt"
)

type Tag struct { string }

var VOID_TYPE = &Tag{"VOID"}
var BOOL_TYPE = &Tag{"BOOL"}
var INT_TYPE = &Tag{"INT"}


type DataType struct {
	BaseType *Tag
	SubTypes []*DataValue
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
	Name string
	Symbols map[string]*Symbol
	Parent *SymbolTable
}

func ValueString(dv *DataValue) string {
	//FIXME
	return "???"
}

func TypeString(dt *DataType) string {
	ret := dt.BaseType.string
	if dt.SubTypes != nil {
		ret += "("
		for i, st := range dt.SubTypes {
			if i > 0 {
				ret += ", "
			}
			ret += ValueString(st)
		}
		ret += ")"
	}

	return ret
}

func EmitSymbolTable(st *SymbolTable) {

	for st != nil {
		fmt.Printf("----%s---\n", st.Name)
		for k, v := range st.Symbols {
			fmt.Printf("%v %v = %v\n",
				k,
				TypeString(v.Type),
				ValueString(v.InitialValue))
		}
		st = st.Parent
	}
}

func ResolveGlobals(file_set *FileSet) {
	//FIXME implement
}

var boolType = &DataType{BaseType: BOOL_TYPE}

var trueValue = &DataValue{ Type: boolType }
var falseValue = &DataValue{ Type: boolType }


//FIXME implement
var PredefinedSymbols = &SymbolTable {
	Name: "PREDEFINED",
	Symbols: map[string]*Symbol {
		"True": &Symbol {"True", boolType, trueValue } ,
		"False": &Symbol {"False", boolType, falseValue } },
	Parent: nil }

