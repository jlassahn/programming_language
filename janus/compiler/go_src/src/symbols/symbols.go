
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
	String() string
	Lookup(string) Symbol
	LookupOperator(string) FunctionChoiceSymbol
	Emit(emitParent bool) //FIXME do we really use recursive emits?

	AddConst(name string, dtype DataType, val DataValue) error
	AddVar(name string, dtype DataType) (Symbol, error)

	AddOperator(name string, dtype DataType, isConst bool, impl DataValue) error

}

type symbolTable struct {
	Name string
	Symbols map[string]Symbol
	Operators map[string] FunctionChoiceSymbol
	Parent *symbolTable
}

func NewSymbolTable(name string, parent SymbolTable) *symbolTable {

	var par *symbolTable
	if parent != nil {
		par = parent.(*symbolTable)
	}

	return &symbolTable {
		Name: name,
		Symbols: map[string]Symbol {},
		Operators: map[string]FunctionChoiceSymbol {},
		Parent: par,
	}
}

func (self *symbolTable) String() string { return self.Name }
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

func (self *symbolTable) AddOperator(name string, dtype DataType,
	isConst bool, impl DataValue) error {

	if self.Operators[name] == nil {
		self.Operators[name] = &functionChoiceSymbol {name, nil}
	}

	choices := self.Operators[name]

	return choices.Add(
		&baseSymbol {
			name,
			dtype,
			impl,
			isConst,
			nil,
		})

}

//FIXME should this be part of the SymbolTable interface?
func (self *symbolTable) AddFunction(
	name string, dtype DataType,
	isConst bool, impl DataValue) error {

	if self.Symbols[name] == nil {
		self.Symbols[name] = &functionChoiceSymbol {name, nil}
	}

	choices, ok := self.Symbols[name].(FunctionChoiceSymbol)
	if !ok {
		return fmt.Errorf("multiple symbol definitions for %v", name)
	}

	return choices.Add(
		&baseSymbol {
			name,
			dtype,
			impl,
			isConst,
			nil,
		})

}

func (st *symbolTable) Emit(emitParent bool) {

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

			v.InitialValue().(NamespaceDataValue).AsSymbolTable().Emit(false)
		}

		st = st.Parent
		if !emitParent {
			break
		}
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
func (self *functionChoiceSymbol) IsConst() bool { return true; }
func (self *functionChoiceSymbol) Choices() []Symbol { return self.choices; }
func (self *functionChoiceSymbol) SetGenVal(val interface{}) { }
func (self *functionChoiceSymbol) GetGenVal() interface{} { return nil }

func (self *functionChoiceSymbol) InitialValue() DataValue {
	return &functionChoiceDV { self }
}

func (self *functionChoiceSymbol) Add(x Symbol) error {
	//FIXME
	// if FunctionParamsAmbiguous(params, choices) ...
	// FIXME how to handle error messsages, etc

	for i, old := range self.choices {
		if TypeMatches(old.Type(), x.Type()) {
			if x.InitialValue() != nil {
				if old.InitialValue() != nil {
					return fmt.Errorf("multiple definitions")
				}
				self.choices[i] = x
				return nil
			}
			return nil
		}
	}

	self.choices = append(self.choices, x)
	return nil
}


