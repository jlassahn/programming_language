
package generator

import (
	"output"
)

type GeneratedFile interface {
	OutFile() output.ObjectFile
	EmitComment(msg string, args ...interface{})
	Emit(msg string, args ...interface{})

	MakeResult() *result
	SetMain(mainName string)
	GetMain() string
}

type generatedFile struct {
	outFile output.ObjectFile
	nextID int
	mainName string
}

func NewGeneratedFile(outFile output.ObjectFile) GeneratedFile {

	return &generatedFile {
		outFile: outFile,
		nextID: 0,
		mainName: "",
	}
}

func (self *generatedFile) OutFile() output.ObjectFile {
	return self.outFile
}

func (self *generatedFile) EmitComment(msg string, args ...interface{}) {
	self.outFile.EmitComment(msg, args...)
}

func (self *generatedFile) Emit(msg string, args ...interface{}) {
	self.outFile.Emit(msg, args...)
}

func (self *generatedFile) MakeResult() *result {

	ret := &result { }
	ret.id = self.nextID
	self.nextID ++

	return ret
}

func (self *generatedFile) SetMain(mainName string) {

	self.mainName = mainName
}

func (self *generatedFile) GetMain() string {
	return self.mainName
}

