
package generator

import (
	"fmt"

	"symbols"
)


type Result interface {
	ID() int
	Name() string
	Type() symbols.DataType

	//FIXME possible values include
	// constant
	// temporary result by value
	// variable by reference
	// function choice possibly with base object

	IsConst() bool
	ConstVal() symbols.DataValue

	IsFunctionChoice() bool
	FunctionChoice() symbols.FunctionChoiceSymbol

	IsVariableRef() bool
	IsGlobalRef() bool

	String() string
	LLVMType() string
	LLVMVal() string
}

type result struct {
	id int
	name string
	dtype symbols.DataType
	constVal symbols.DataValue
	functionChoice symbols.FunctionChoiceSymbol
	isVariableRef bool
	isGlobalRef bool
}

func NewTempVal(fp GeneratedFile, dtype symbols.DataType) Result {

	ret := fp.MakeResult()
	ret.dtype = dtype
	ret.name = "tmp"
	ret.constVal = nil

	return ret
}

func NewNamedVal(fp GeneratedFile, dtype symbols.DataType, name string) Result {

	ret := fp.MakeResult()
	ret.dtype = dtype
	ret.name = name
	ret.constVal = nil
	ret.isVariableRef = true
	ret.isGlobalRef = false

	return ret
}

func NewGlobalVal(fp GeneratedFile,
	dtype symbols.DataType, name string) Result {

	ret := &result {}
	ret.dtype = dtype
	ret.name = name
	ret.isGlobalRef = true

	return ret
}

func NewDataVal(dval symbols.DataValue) Result {

	ret := &result { }
	ret.dtype = dval.Type()
	ret.name = "INVALID"
	ret.constVal = dval

	return ret
}

func NewTypedDataVal(dtype symbols.DataType, dval symbols.DataValue) Result {

	ret := &result { }
	ret.dtype = dtype
	ret.name = "INVALID"
	ret.constVal = dval

	return ret
}

func NewFunctionChoiceResult(fn symbols.FunctionChoiceSymbol) Result {

	ret := &result { }
	ret.dtype = fn.Type()
	ret.name = fn.Name()
	ret.functionChoice = fn

	return ret
}


func (self *result) String() string { return self.LLVMVal() }
func (self *result) ID() int { return self.id }
func (self *result) Name() string { return self.name }
func (self *result) Type() symbols.DataType { return self.dtype }
func (self *result) IsConst() bool { return self.constVal != nil }
func (self *result) ConstVal() symbols.DataValue { return self.constVal }
func (self *result) IsVariableRef() bool { return self.isVariableRef }
func (self *result) IsGlobalRef() bool { return self.isGlobalRef }

func (self *result) IsFunctionChoice() bool {
	return self.functionChoice != nil
}

func (self *result) FunctionChoice() symbols.FunctionChoiceSymbol {
	return self.functionChoice
}

func (self *result) LLVMType() string {
	return MakeLLVMType(self.dtype)
}

func (self *result) LLVMVal() string {
	if self.IsConst() {
		return MakeLLVMConst(self.constVal)
	} else if self.IsFunctionChoice() {
		return "UNRESOLVED_FUNCTION_CHOICE"
	} else if self.IsGlobalRef() {
		return fmt.Sprintf("@%s", self.name)
	} else{
		return fmt.Sprintf("%%%s_%d", self.name, self.id)
	}
}

