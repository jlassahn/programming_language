
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

	Add(x Symbol) error
}

type SymbolTable interface {
	Lookup(string) Symbol
	LookupOperator(string) FunctionChoiceSymbol
	Emit()

	AddConst(name string, dtype DataType, val DataValue) error

	AddOperator(name string, retType DataType, params []FunctionParameter,
		isConst bool, impl DataValue) error
}

type symbolTable struct {
	Name string
	Symbols map[string]Symbol
	Operators map[string] FunctionChoiceSymbol
	Parent *symbolTable
}

func NewSymbolTable(name string, parent *symbolTable) *symbolTable {

	return &symbolTable {
		Name: name,
		Symbols: map[string]Symbol {},
		Operators: map[string]FunctionChoiceSymbol {},
		Parent: parent,
	}
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

func (self *symbolTable) AddConst(
	name string, dtype DataType, val DataValue) error {

	if self.Symbols[name] != nil {
		return fmt.Errorf("redefinition of symbol %v", name)
	}

	self.Symbols[name] = &baseSymbol { name, dtype, val, true }
	return nil
}

func (self *symbolTable) AddOperator(
	name string, retType DataType, params []FunctionParameter,
	isConst bool, impl DataValue) error {

	if self.Operators[name] == nil {
		self.Operators[name] = &functionChoiceSymbol {name, nil}
	}

	choices := self.Operators[name]

	return choices.Add(
		&baseSymbol {
			name,
			&functionDT {retType, params, false},
			impl,
			isConst,
		})

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

func (self *baseSymbol) String() string {
	return self.name + ":" + self.dtype.String()
}


type functionChoiceSymbol struct {
	name string
	choices []Symbol
}

func (self *functionChoiceSymbol) Name() string { return self.name; }
func (self *functionChoiceSymbol) Type() DataType { return FunctionChoiceType; }
func (self *functionChoiceSymbol) InitialValue() DataValue { return nil; }
func (self *functionChoiceSymbol) IsConst() bool { return true; }
func (self *functionChoiceSymbol) Choices() []Symbol { return self.choices; }

func (self *functionChoiceSymbol) Add(x Symbol) error {
	//FIXME
	// if FunctionParamsAmbiguous(params, choices) ...
	// FIXME how to handle error messsages, etc

	self.choices = append(self.choices, x)
	return nil
}



var predefinedSymbols *symbolTable;

func PredefinedSymbols() *symbolTable {
	if predefinedSymbols == nil {
		predefinedSymbols = buildPredefinedSymbols()
	}
	return predefinedSymbols
}


//FIXME reorganize and correct
func buildPredefinedSymbols() *symbolTable {

	syms := NewSymbolTable("PREDEFINED", nil)

	syms.AddConst("True", BoolType, TrueValue)
	syms.AddConst("False", BoolType, FalseValue)

	syms.AddOperator("+", Real64Type, []FunctionParameter {
		{"a", Real64Type, false},
		{"b", Real64Type, true},
	},
	true, nil) //FIXME nil implementation

	syms.AddOperator("+", Int64Type, []FunctionParameter {
		{"a", Int64Type, false},
		{"b", Int64Type, true},
	},
	true, IntrinsicAddInt64)

	syms.AddOperator("/", Real64Type, []FunctionParameter {
		{"a", Real64Type, false},
		{"b", Real64Type, true},
	},
	true, IntrinsicDivReal64)

	//FIXME should  be Int64, >Real64
	//FIXME should have other IntXXX versions
	syms.AddOperator("/", Real64Type, []FunctionParameter {
		{"a", Int64Type, false},
		{"b", Int64Type, true},
	},
	true, IntrinsicDivReal64)

	return syms
}

