
package generator

import (
	"fmt"

	"symbols"
)

type GeneratedFunction interface {
	File() GeneratedFile
	ReturnType() symbols.DataType

	Emit()
	SetReturnType(dtype symbols.DataType)
	AddParameter(name string, dtype symbols.DataType) Result
	AddPrologue(x string, args ...interface{})
	AddBody(x string, args ...interface{})
}

type generatedFunction struct {
	fp GeneratedFile
	name string
	returnType symbols.DataType
	params []Result
	prologue []string
	body []string
}

func NewGeneratedFunction(fp GeneratedFile, name string) GeneratedFunction {

	return &generatedFunction {
		fp: fp,
		name: name,
		returnType: symbols.VoidType,
		params: nil,
	}
}

func (self *generatedFunction) File() GeneratedFile {
	return self.fp
}

func (self *generatedFunction) ReturnType() symbols.DataType {
	return self.returnType
}

func (self *generatedFunction) Emit() {

	self.fp.Emit("define %v @%v(", MakeLLVMType(self.returnType), self.name)
	for i,param := range self.params {
		term := ","
		if i==len(self.params)-1 {
			term = ""
		}
		self.fp.Emit("\t%v%v", param.LLVMType(), term)
	}
	self.fp.Emit(") {")

	for _,x := range self.prologue {
		self.fp.Emit("%s", x)
	}

	self.fp.Emit("")

	for _,x := range self.body {
		self.fp.Emit("%s", x)
	}

	self.fp.Emit("")
	//FIXME fake, only accessible if the function doesn't call return on
	//      some path
	if self.returnType == symbols.VoidType {
		self.fp.Emit("\tret void")
	} else {
		self.fp.Emit("\tret %v zeroinitializer", MakeLLVMType(self.returnType))
	}
	self.fp.Emit("}")
}

func (self *generatedFunction) SetReturnType(dtype symbols.DataType) {
	self.returnType = dtype
}

func (self *generatedFunction) AddParameter(name string, dtype symbols.DataType) Result {

	ret := NewLocalVal(self.fp, dtype, name)
	self.params = append(self.params, ret)

	return ret
}

func (self *generatedFunction) AddPrologue(x string, args ...interface{}) {
	s := fmt.Sprintf(x, args...)
	self.prologue = append(self.prologue, s)
}

func (self *generatedFunction) AddBody(x string, args ...interface{}) {
	s := fmt.Sprintf(x, args...)
	self.body = append(self.body, s)
}

