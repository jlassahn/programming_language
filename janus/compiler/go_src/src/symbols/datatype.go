
package symbols

import (
)

type Tag struct { string }
func (self *Tag) String() string {
	return self.string
}

type DTypeParameter struct {
	Number int64
	DType DataType
}

type DataType interface {
	String() string
	Base() *Tag
	SubTypes() []DTypeParameter
	Members() map[string]Symbol
}

type FunctionParameter struct {
	Name string
	DType DataType
	AutoConvert bool
}

func (self FunctionParameter) String() string {

	ret := self.Name + " "
	if self.AutoConvert {
		ret = ret + ">"
	}
	ret = ret + self.DType.String()

	return ret
}

type FunctionDataType interface {
	DataType
	ReturnType() DataType
	Parameters() []FunctionParameter
	IsMethod() bool

	AddParam(name string, dtype DataType, auto bool) FunctionDataType
}

type functionDT struct {
	returnType DataType
	parameters []FunctionParameter
	isMethod bool
}

func NewFunction(retType DataType) FunctionDataType {
	return &functionDT {
		returnType: retType,
		parameters: nil,
		isMethod: false,
	}
}

func (self *functionDT) AddParam(
	name string, dtype DataType, auto bool) FunctionDataType {

	param := FunctionParameter {
		Name: name,
		DType: dtype,
		AutoConvert: auto,
	}

	self.parameters = append(self.parameters, param)
	return self
}

func (self *functionDT) String() string {

	ret := ""
	if self.isMethod {
		ret ="METHOD("
	} else {
		ret ="FUNCTION("
	}
	for i,x := range self.parameters {
		if i > 0 {
			ret = ret + ", "
		}
		ret = ret + x.String()
	}
	ret = ret + ")->"
	ret = ret + self.returnType.String()

	return ret
}

func (self *functionDT) Base() *Tag {
	return FUNCTION_TYPE
}

func (self *functionDT) SubTypes() []DTypeParameter {
	return nil
}

func (self *functionDT) Members() map[string]Symbol {
	return nil
}

func (self *functionDT) ReturnType() DataType {
	return self.returnType
}

func (self *functionDT) Parameters() []FunctionParameter {
	return self.parameters
}

func (self *functionDT) IsMethod() bool {
	return self.isMethod
}

type simpleDT struct {
	tag *Tag
	members map[string]Symbol
}

func (self *simpleDT) Base() *Tag {
	return self.tag
}

func (self *simpleDT) SubTypes() []DTypeParameter {
	return nil
}

func (self *simpleDT) Members() map[string]Symbol {
	return self.members
}

func (self *simpleDT) String() string {
	return self.tag.string
}

