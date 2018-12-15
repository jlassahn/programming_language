
package symbols

import (
	"fmt"

	"output"
)

//FIXME better ways to look up polymorphic functions

type Symbol interface {
	Name() string
	Type() DataType
	InitialValue() DataValue
	IsConst() bool
	SetGenVal(val interface{})
	GetGenVal() interface{}
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
	AddVar(name string, dtype DataType) (Symbol, error)

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

//FIXME return (Symbol, error)
func (self *symbolTable) AddConst(
	name string, dtype DataType, val DataValue) error {

	if self.Symbols[name] != nil {
		return fmt.Errorf("redefinition of symbol %v", name)
	}

	self.Symbols[name] = &baseSymbol { name, dtype, val, true, nil }
	return nil
}

func (self *symbolTable) AddVar(name string, dtype DataType) (Symbol, error) {

	if self.Symbols[name] != nil {
		return nil, fmt.Errorf("redefinition of symbol %v", name)
	}

	sym := &baseSymbol { name, dtype, nil, false, nil }
	self.Symbols[name] = sym
	return sym, nil
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
			nil,
		})

}

func (st *symbolTable) Emit() {

	for st != nil {
		output.Emit("  Symbol Table: %v", st.Name)
		output.Emit("    Symbols:")
		for _,k := range SortedKeys(st.Symbols) {
			v := st.Symbols[k]
			if v.Type() == FunctionChoiceType {
				output.Emit("      %v %v",
					k,
					v.Type())
				for _, op := range v.(FunctionChoiceSymbol).Choices() {
					output.Emit("        %v = %v",
						op.Type(),
						op.InitialValue())
				}
			} else {
				output.Emit("      %v %v = %v",
					k,
					v.Type(),
					v.InitialValue())
			}
		}

		output.Emit("    Operators:")
		for k, v := range st.Operators {
			output.Emit("      %v %v",
				k,
				v.Type())
			for _, op := range v.Choices() {
				output.Emit("        %v = %v",
					op.Type(),
					op.InitialValue())
			}
		}
		output.Emit("")

		for _,v := range st.Symbols {
			if v.Type() != NamespaceType {
				continue
			}

			v.InitialValue().(NamespaceDataValue).AsSymbolTable().Emit()
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
	genVal interface{}
}

func (self *baseSymbol) Name() string { return self.name }
func (self *baseSymbol) Type() DataType { return self.dtype }
func (self *baseSymbol) InitialValue() DataValue { return self.initialValue }
func (self *baseSymbol) IsConst() bool { return self.isConst }
func (self *baseSymbol) SetGenVal(val interface{}) { self.genVal = val }
func (self *baseSymbol) GetGenVal() interface{} { return self.genVal }

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
func (self *functionChoiceSymbol) SetGenVal(val interface{}) { }
func (self *functionChoiceSymbol) GetGenVal() interface{} { return nil }

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

	syms.Symbols["__system"] = buildInternalSymbols()

	syms.AddConst("TRUE", BoolType, TrueValue)
	syms.AddConst("FALSE", BoolType, FalseValue)

	syms.AddConst("CType", CTypeType, &typeDV{CTypeType, CTypeType})
	syms.AddConst("Bool", CTypeType, &typeDV{CTypeType, BoolType})
	syms.AddConst("Int32", CTypeType, &typeDV{CTypeType, Int32Type})

	syms.AddOperator("+", Real64Type, []FunctionParameter {
		{"a", Real64Type, false},
		{"b", Real64Type, true},
	},
	true, nil) //FIXME nil implementation

	syms.AddOperator("+", Int64Type, []FunctionParameter {
		{"a", Int64Type, false},
		{"b", Int64Type, true},
	},
	true, &intrinsicDV{"add_Int64"} )

	syms.AddOperator("+", Int32Type, []FunctionParameter {
		{"a", Int32Type, false},
		{"b", Int32Type, true},
	},
	true, &intrinsicDV{"add_Int64"})

	syms.AddOperator("/", Real64Type, []FunctionParameter {
		{"a", Real64Type, false},
		{"b", Real64Type, true},
	},
	true, &intrinsicDV{"div_Real64"})

	//FIXME should  be Int64, >Real64
	//FIXME should have other IntXXX versions
	syms.AddOperator("/", Real64Type, []FunctionParameter {
		{"a", Int64Type, false},
		{"b", Int64Type, true},
	},
	true, &intrinsicDV{"div_Real64"})

	return syms
}

func buildInternalSymbols() *baseSymbol {

	name := "PREDEFINED:__system"
	newTable := NewSymbolTable(name, nil)

	val := &namespaceDV {
		value: newTable,
	}

	return &baseSymbol {
		name: name,
		dtype: NamespaceType,
		initialValue: val,
		isConst: true,
	}
}

