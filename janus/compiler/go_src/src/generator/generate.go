
package generator

import (
	"fmt"

	"output"
	"parser"
	"symbols"
)

//FIXME add methods to Generated... so we don't have to keep saying
// MakeLLVMType(self.Result().Type())

type GeneratedTag interface {
	ID() int
	Type() symbols.DataType
	String() string
}

type GeneratedStatement interface {
	Result() GeneratedTag
	IsConst() bool
	ConstVal() symbols.DataValue
	String() string
}

type GeneratedFunction interface {
	Emit(file GeneratedFile)
	SetReturnType(dtype symbols.DataType)
	AddParameter(name string, dtype symbols.DataType) GeneratedStatement
	AddPrologue(x string, args ...interface{})
	AddBody(x string, args ...interface{})
}

type GeneratedGlobalDef interface {
}

type GeneratedFile interface {
	OutFile() output.ObjectFile
	EmitComment(msg string, args ...interface{})
	Emit(msg string, args ...interface{})
	MakeTag(tagOut *generatedTag, dtype symbols.DataType, name string)
}


type generatedTag struct {
	id int
	name string
	dtype symbols.DataType
}

func (self *generatedTag) ID() int {
	return self.id
}

func (self *generatedTag) Type() symbols.DataType {
	return self.dtype
}

func (self *generatedTag) String() string {
	return fmt.Sprintf("%%%s_%d", self.name, self.id)
}

type generatedStatement struct {
	tag generatedTag
	isConst bool
	constVal symbols.DataValue
}

func NewTempVal(fp GeneratedFile, dtype symbols.DataType) *generatedStatement {

	ret := &generatedStatement { }
	fp.MakeTag(&ret.tag, dtype, "tmp")

	return ret
}

func NewNamedVal(fp GeneratedFile, dtype symbols.DataType, name string) *generatedStatement {

	ret := &generatedStatement { }
	fp.MakeTag(&ret.tag, dtype, name)

	return ret
}

func (self *generatedStatement) String() string {
	if self.isConst {
		return self.constVal.String()
	} else {
		return self.tag.String()
	}
}
func (self *generatedStatement) Result() GeneratedTag { return &self.tag }
func (self *generatedStatement) IsConst() bool { return self.isConst }
func (self *generatedStatement) ConstVal() symbols.DataValue { return  self.constVal }

type generatedParam struct {
	tag generatedTag
	index int
}

func NewGeneratedParam(fp GeneratedFile, i int, dtype symbols.DataType,
	name string) GeneratedStatement {

	ret := &generatedParam { }
	fp.MakeTag(&ret.tag, dtype, name)
	ret.index = i

	return ret
}

func (self *generatedParam) String() string {
	return fmt.Sprintf("parameter(%v)", self.tag.name)
}

func (self *generatedParam) Result() GeneratedTag {
	return &self.tag
}

func (self *generatedParam) IsConst() bool { return false }
func (self *generatedParam) ConstVal() symbols.DataValue { return nil }

type generatedFile struct {
	outFile output.ObjectFile
	nextID int
}

func NewGeneratedFile(outFile output.ObjectFile) GeneratedFile {

	return &generatedFile {
		outFile: outFile,
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

func (self *generatedFile) MakeTag(tagOut *generatedTag, dtype symbols.DataType, name string) {

	tagOut.id = self.nextID
	self.nextID ++
	tagOut.name = name
	tagOut.dtype = dtype
}

type generatedFunction struct {
	fp GeneratedFile
	name string
	returnType symbols.DataType
	params []GeneratedStatement
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

func (self *generatedFunction) Emit(fp GeneratedFile) {

	fp.Emit("define %v @%v(", MakeLLVMType(self.returnType), self.name)
	for i,param := range self.params {
		term := ","
		if i==len(self.params)-1 {
			term = ""
		}
		fp.Emit("\t%v%v", MakeLLVMType(param.Result().Type()), term)
	}
	fp.Emit(") {")

	for _,x := range self.prologue {
		fp.Emit("%s", x)
	}

	fp.Emit("")

	for _,x := range self.body {
		fp.Emit("%s", x)
	}

	fp.Emit("")
	//FIXME fake, only accessible if the function doesn't call return on
	//      some path
	fp.Emit("endpoint:")
	if self.returnType == symbols.VoidType {
		fp.Emit("\tret void")
	} else {
		fp.Emit("\tret %v zeroinitializer", MakeLLVMType(self.returnType))
	}
	fp.Emit("}")
}

func (self *generatedFunction) SetReturnType(dtype symbols.DataType) {
	self.returnType = dtype
}

func (self *generatedFunction) AddParameter(name string, dtype symbols.DataType) GeneratedStatement {

	ret := NewNamedVal(self.fp, dtype, name)
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


func GenerateCode(fileSet *symbols.FileSet, objFile output.ObjectFile) {

	genFile := NewGeneratedFile(objFile)

	mods := fileSet.RootModule.GetModuleList()

	for _,mod := range mods {
		GenerateVariables(fileSet, genFile, mod)
	}

	for _,mod := range mods {
		GenerateFunctions(fileSet, genFile, mod)
	}
}

func GenerateVariables(fileSet *symbols.FileSet, fp GeneratedFile, mod *symbols.Module) {
	//FIXME implement
	fp.EmitComment("")
	fp.EmitComment("generating variables for %v", mod.Name)
	fp.EmitComment("")
}

func GenerateFunctions(fileSet *symbols.FileSet, fp GeneratedFile, mod *symbols.Module) {

	fp.EmitComment("")
	fp.EmitComment("generating functions for %v", mod.Name)
	fp.EmitComment("")

	for _,name := range symbols.SortedKeys(mod.LocalSymbols.Symbols) {
		sym := mod.LocalSymbols.Symbols[name]
		choice, ok := sym.(symbols.FunctionChoiceSymbol)
		if !ok {
			continue
		}

		for _, fn := range choice.Choices() {
			GenerateFunction(fp, mod, fn)
		}
	}

	//mod.LocalSymbols.Operators
}

func GenerateFunction(
	fp GeneratedFile,
	mod *symbols.Module,
	fn symbols.Symbol) {

	el := fn.InitialValue().(symbols.CodeDataValue).AsParseElement()
	file := fn.InitialValue().(symbols.CodeDataValue).AsSourceFile()
	dtype := fn.Type().(symbols.FunctionDataType)
	name := MakeSymbolName(file.Options.ModuleName, dtype, fn.Name())

	symbolTable := symbols.NewSymbolTable(
		fmt.Sprintf("local@%d", el.FilePos().Line),
		file.FileSymbols)

	ctx := &symbols.EvalContext {
		Symbols: symbolTable,
	}

	//FIXME do a first pass to find labels?

	genFunc := NewGeneratedFunction(fp, name)

	for i,param := range dtype.Parameters() {
		genParam := genFunc.AddParameter(param.Name, param.DType)
		sym, err := ctx.Symbols.AddVar(param.Name, param.DType)
		if err != nil {
			parser.Error(el.FilePos(), "%v", err)
			return 
		}
		sym.SetGenVal(genParam)

		genFunc.AddPrologue("\t%v = alloca %v", genParam.String(),
			MakeLLVMType(genParam.Result().Type()))
		genFunc.AddBody("\tstore %v %%%d, %v* %v",
			MakeLLVMType(genParam.Result().Type()), i,
			MakeLLVMType(genParam.Result().Type()),
			genParam.Result().String())
	}
	genFunc.SetReturnType(dtype.ReturnType())

	//FIXME
	ctx.Symbols.Emit()

	for _,elem := range el.Children() {
		GenerateStatement(fp, genFunc, ctx, elem)
	}

	output.FIXMEDebug("generating %v %v %v", dtype, el, file)
	output.FIXMEDebug("name %v", name)

	genFunc.Emit(fp)
}


func GenerateStatement(fp GeneratedFile, genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) {

	GenerateExpression(fp, genFunc, ctx, el)
}


func MakeSymbolName(mod []string, dtype symbols.DataType, name string) string {

	ret := ""
	for _,m := range mod {
		ret = ret + m + "."
	}

	ret = ret + name

	if dtype.Base() == symbols.FUNCTION_TYPE {
		ret = "f_" + ret + MakeTypeName(dtype)
	} else if dtype.Base() == symbols.LABEL_TYPE {
		ret = "l_" + ret
	} else {
		ret = "d_" + ret
	}

	return ret
}

func MakeTypeName(dtype symbols.DataType) string {

	ret := "-" + dtype.Base().String()

	if dtype.Base() == symbols.FUNCTION_TYPE {
		ftype := dtype.(symbols.FunctionDataType)
		for _,param := range ftype.Parameters() {
			ret = ret + MakeTypeName(param.DType)
		}
		ret = ret + MakeTypeName(ftype.ReturnType())
	} else {
		for _,param := range dtype.SubTypes() {
			if param.DType == nil {
				ret = ret + "-" + string(param.Number) + "$"
			} else {
				ret = ret + MakeTypeName(param.DType)
			}
		}
	}

	ret = ret + "$"
	return ret
}

func MakeLLVMType(dtype symbols.DataType) string {
	return "i32" //FIXME fake
}

