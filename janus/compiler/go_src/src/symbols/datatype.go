
package symbols

import (
)

type Tag struct { string }

type DataType interface {
	String() string
	Base() *Tag
	SubTypes() []DataValue
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
}

type functionDT struct {
	returnType DataType
	parameters []FunctionParameter
	isMethod bool
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

func (self *functionDT) SubTypes() []DataValue {
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
	Tag
}

func (self *simpleDT) Base() *Tag {
	return &self.Tag
}

func (self *simpleDT) SubTypes() []DataValue {
	return nil
}

func (self *simpleDT) String() string {
	return self.string
}

