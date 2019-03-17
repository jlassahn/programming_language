
package generator

import (
	"fmt"

	"output"
	"symbols"
)

type ResultTag struct { string }
func (self *ResultTag) String() string {
	return self.string
}

var CONST_RESULT = &ResultTag{"CONST_RESULT"}
var TEMP_RESULT = &ResultTag{"TEMP_RESULT"}
var LOCAL_RESULT = &ResultTag{"LOCAL_RESULT"}
var GLOBAL_RESULT = &ResultTag{"GLOBAL_RESULT"}
var FN_CHOICE_RESULT = &ResultTag{"FN_CHOICE_RESULT"}

type Result interface {
	ID() int
	Name() string
	Type() symbols.DataType //FIXME DType()?
	Tag() *ResultTag
	
	//FIXME possible values include
	// constant
	// temporary result by value
	// variable by reference
	// function choice possibly with base object

	IsConst() bool
	ConstVal() symbols.DataValue

	IsFunctionChoice() bool
	FunctionChoice() symbols.FunctionChoiceSymbol
	HasBaseObject() bool
	BaseObject() Result

	IsVariableRef() bool

	String() string
	LLVMType() string
	LLVMVal() string
}

type result struct {
	tag *ResultTag
	id int
	name string
	dtype symbols.DataType
	constVal symbols.DataValue
	functionChoice symbols.FunctionChoiceSymbol
	baseObject Result
}

func NewTempVal(fp GeneratedFile, dtype symbols.DataType) Result {

	ret := fp.MakeResult()
	ret.tag = TEMP_RESULT
	ret.dtype = dtype
	ret.name = "tmp"
	ret.constVal = nil

	return ret
}

func NewLocalVal(
	fp GeneratedFile,
	dtype symbols.DataType,
	name string,
) Result {

	ret := fp.MakeResult()
	ret.tag = LOCAL_RESULT
	ret.dtype = dtype
	ret.name = name
	ret.constVal = nil

	return ret
}

func NewGlobalVal(fp GeneratedFile,
	dtype symbols.DataType, name string) Result {

	ret := &result {}
	ret.tag = GLOBAL_RESULT
	ret.dtype = dtype
	ret.name = name

	return ret
}

func NewConstVal(dval symbols.DataValue) Result {

	ret := &result { }
	ret.tag = CONST_RESULT
	ret.dtype = dval.Type()
	ret.name = "INVALID"
	ret.constVal = dval

	return ret
}

func NewZeroVal(dtype symbols.DataType) Result {

	ret := &result { }
	ret.tag = CONST_RESULT
	ret.dtype = dtype
	ret.name = "INVALID"
	ret.constVal = nil

	return ret
}

func NewFunctionChoiceResult(fn symbols.FunctionChoiceSymbol) Result {

	ret := &result { }
	ret.tag = FN_CHOICE_RESULT
	ret.dtype = fn.Type()
	ret.name = fn.Name()
	ret.functionChoice = fn

	return ret
}

func NewMethodChoiceResult(fn symbols.FunctionChoiceSymbol, base Result) Result {

	ret := &result { }
	ret.tag = FN_CHOICE_RESULT
	ret.dtype = fn.Type()
	ret.name = fn.Name()
	ret.functionChoice = fn
	ret.baseObject = base

	return ret
}

func (self *result) String() string { return self.LLVMVal() }
func (self *result) Tag() *ResultTag { return self.tag }
func (self *result) ID() int { return self.id }
func (self *result) Name() string { return self.name }
func (self *result) Type() symbols.DataType { return self.dtype }
func (self *result) IsConst() bool { return self.tag == CONST_RESULT }
func (self *result) ConstVal() symbols.DataValue { return self.constVal }

func (self *result) IsVariableRef() bool {
	return self.tag == LOCAL_RESULT || self.tag == GLOBAL_RESULT
}

func (self *result) IsFunctionChoice() bool {
	return self.tag == FN_CHOICE_RESULT
}

func (self *result) FunctionChoice() symbols.FunctionChoiceSymbol {
	return self.functionChoice
}

func (self *result) HasBaseObject() bool {
	return self.baseObject != nil
}

func (self *result) BaseObject() Result {
	return self.baseObject
}

func (self *result) LLVMType() string {
	return MakeLLVMType(self.dtype)
}

func (self *result) LLVMVal() string {

	switch self.tag {
	case CONST_RESULT:
		if self.constVal == nil {
			return "zeroinitializer"
		}
		return MakeLLVMConst(self.constVal)
	case TEMP_RESULT:
		return fmt.Sprintf("%%%s_%d", self.name, self.id)
	case LOCAL_RESULT:
		return fmt.Sprintf("%%%s_%d", self.name, self.id)
	case GLOBAL_RESULT:
		return fmt.Sprintf("@%s", self.name)
	case FN_CHOICE_RESULT:
		return "UNRESOLVED_FUNCTION_CHOICE"

	default:
		output.FatalError("invalid result tag")
		return "ERROR"
	}
}

