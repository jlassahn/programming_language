
package symbols

import (
	"fmt"
)

//FIXME better ways to look up polymorphic functions

type Symbol interface {
	Name() string
	Type() DataType
	InitialValue() DataValue
	IsConst() bool
}

type FunctionChoiceSymbol interface {
	Symbol
	Choices() []Symbol
}

type SymbolTable interface {
	Lookup(string) Symbol
	LookupOperator(string) FunctionChoiceSymbol
	Emit()
}

type symbolTable struct {
	Name string
	Symbols map[string]Symbol
	Operators map[string]FunctionChoiceSymbol
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

func (self *symbolTable) LookupOperator(x string) FunctionChoiceSymbol {

	ret := self.Operators[x]
	if ret != nil {
		return ret
	}
	if self.Parent == nil {
		return nil
	}

	return self.Parent.LookupOperator(x)
}

func (st *symbolTable) Emit() {

	for st != nil {
		fmt.Printf("----%s---\n", st.Name)
		fmt.Printf("Symbols:\n")
		for k, v := range st.Symbols {
			fmt.Printf("  %v %v = %v\n",
				k,
				v.Type(),
				v.InitialValue())
		}

		//FIXME more detail about operators
		fmt.Printf("Operators:\n")
		for k, v := range st.Operators {
			fmt.Printf("  %v %v\n",
				k,
				v.Type())
			for _, op := range v.Choices() {
				fmt.Printf("    %v = %v\n",
					op.Type(),
					op.InitialValue())
			}
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
	isConst bool
}

func (self *baseSymbol) Name() string { return self.name; }
func (self *baseSymbol) Type() DataType { return self.dtype; }
func (self *baseSymbol) InitialValue() DataValue { return self.initialValue; }
func (self *baseSymbol) IsConst() bool { return self.isConst; }


type functionChoiceSymbol struct {
	name string
	choices []Symbol
}

func (self *functionChoiceSymbol) Name() string { return self.name; }
func (self *functionChoiceSymbol) Type() DataType { return FunctionChoiceType; }
func (self *functionChoiceSymbol) InitialValue() DataValue { return nil; }
func (self *functionChoiceSymbol) IsConst() bool { return true; }
func (self *functionChoiceSymbol) Choices() []Symbol { return self.choices; }


//FIXME reorganize and correct

var add_op = []Symbol {
		&baseSymbol { "+",
			&functionDT {
			Int64Type, []FunctionParameter{
				{"a", Int64Type, false},
				{"b", Int64Type, true},
			}, false },
		nil,
		true },
}

/*
var add_op = &functionchoiceDT {
	choices: []FunctionDataType {
		&functionDT{
			Int64Type, []FunctionParameter{
				{"a", Int64Type, false},
				{"b", Int64Type, true},
			}, false },
	},
}

var div_op = &functionchoiceDT {
	choices: []FunctionDataType {
		&functionDT{
			Real64Type, []FunctionParameter{
				{"a", Int64Type, true},
				{"b", Int64Type, true},
			}, false },
	},
}
*/

func PredefinedSymbols() *symbolTable {
	if predefinedSymbols == nil {
		predefinedSymbols = buildPredefinedSymbols()
	}
	return predefinedSymbols
}


func buildPredefinedSymbols() *symbolTable {
	return nil //FIXME
}

var predefinedSymbols = &symbolTable {
	Name: "PREDEFINED",
	Symbols: map[string]Symbol {
		"True": &baseSymbol {"True", BoolType, TrueValue, true } ,
		"False": &baseSymbol {"False", BoolType, FalseValue, true },
	},
	Operators: map[string]FunctionChoiceSymbol {
		"+": &functionChoiceSymbol {"+", add_op },
		/*
		"/": &baseSymbol {"+", div_op, nil, true },
		*/
	},
	Parent: nil }

