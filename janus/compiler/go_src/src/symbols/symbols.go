
package symbols

import (
	"fmt"
)

type Symbol interface {
	Name() string
	Type() DataType
	InitialValue() DataValue
}

type SymbolTable interface {
	Lookup(string) Symbol
	Emit()
}

type symbolTable struct {
	Name string
	Symbols map[string]Symbol
	Parent *symbolTable
}

func (self *symbolTable) Lookup(x string) Symbol {

	ret := self.Symbols[x]
	if ret != nil {
		return ret
	}
	if self.Parent == nil {
		return nil
	}

	return self.Parent.Lookup(x)
}

func (st *symbolTable) Emit() {

	for st != nil {
		fmt.Printf("----%s---\n", st.Name)
		for k, v := range st.Symbols {
			fmt.Printf("%v %v = %v\n",
				k,
				TypeString(v.Type()),
				v.InitialValue().ValueAsString())
		}
		st = st.Parent
	}
}

func ResolveGlobals(file_set *FileSet) {
	//FIXME implement
}


//FIXME organize

type baseSymbol struct {
	name string
	dtype DataType
	initialValue DataValue
}

func (self *baseSymbol) Name() string { return self.name; }
func (self *baseSymbol) Type() DataType { return self.dtype; }
func (self *baseSymbol) InitialValue() DataValue { return self.initialValue; }

//FIXME implement
var PredefinedSymbols = &symbolTable {
	Name: "PREDEFINED",
	Symbols: map[string]Symbol {
		"True": &baseSymbol {"True", BoolType, TrueValue } ,
		"False": &baseSymbol {"False", BoolType, FalseValue } },
	Parent: nil }

