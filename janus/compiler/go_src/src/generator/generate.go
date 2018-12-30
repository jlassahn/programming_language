
package generator

import (
	"fmt"

	"output"
	"parser"
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

	output.FIXMEDebug("setting main: %v", mainName)
	self.mainName = mainName
}

func (self *generatedFile) GetMain() string {
	return self.mainName
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

	GenerateHeader(genFile)

	for _,mod := range mods {
		GenerateVariables(fileSet, genFile, mod)
	}

	for _,mod := range mods {
		GenerateFunctions(fileSet, genFile, mod)
	}

	mainName := genFile.GetMain()
	if mainName != "" {
		genFile.EmitComment("")
		genFile.EmitComment("main entrypoint")
		genFile.EmitComment("")
		genFile.Emit("@main = alias void(), void()* @%v", mainName)
	}

}

func GenerateHeader(fp GeneratedFile) {
	fp.EmitComment("")
	fp.EmitComment("global declarations")
	fp.EmitComment("")

	//FIXME organize better
	fp.Emit("declare double @llvm.sqrt.f64(double)")

}

func GenerateVariables(fileSet *symbols.FileSet, fp GeneratedFile, mod *symbols.Module) {
	fp.EmitComment("")
	fp.EmitComment("generating variables for %v", mod.Name)
	fp.EmitComment("")

	/*
	//FIXME may need to totally rethink how variable assignments happen,
	//      e.g. we want def fn() = thing; to not generate a second copy
	//      of the code.
	for _,name := range symbols.SortedKeys(mod.LocalSymbols.Symbols) {

		sym := mod.LocalSymbols.Symbols[name]
		output.FIXMEDebug("FIXME generate global %v", sym)
	}
	*/

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

	//FIXME implement
	//mod.LocalSymbols.Operators
}

func GenerateFunction(
	fp GeneratedFile,
	mod *symbols.Module,
	fn symbols.Symbol) {

	fdef := fn.InitialValue()
	if fdef == nil {
		//FIXME what about externally defined functions?
		//output.Warning("no definition for function %v", fn)
		GenerateExternFunction(fp, mod, fn)
		return
	}

	if fdef.Tag() != symbols.CODE_VALUE {
		output.FIXMEDebug("FIXME function %v is %v not code", fn, fdef.Type())
		return
	}

	el := fn.InitialValue().(symbols.CodeDataValue).AsParseElement()
	file := fn.InitialValue().(symbols.CodeDataValue).AsSourceFile()
	dtype := fn.Type().(symbols.FunctionDataType)
	name := MakeSymbolName(mod.Path, dtype, fn.Name())

	if fn.Name() == "Main" {
		if dtype.ReturnType() != symbols.VoidType ||
		len(dtype.Parameters()) != 0 {
			parser.Error(el.FilePos(), "Main must be Main()->Void")
		} else if fp.GetMain() != "" {
			parser.Error(el.FilePos(), "multiple definitions of Main")
		} else {
			fp.SetMain(name)
		}
	}

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

		genFunc.AddPrologue("\t%v = alloca %v", genParam.LLVMVal(),
			genParam.LLVMType())
		genFunc.AddBody("\tstore %v %%%d, %v* %v",
			genParam.LLVMType(), i,
			genParam.LLVMType(),
			genParam.LLVMVal())
	}
	genFunc.SetReturnType(dtype.ReturnType())

	for _,elem := range el.Children() {
		GenerateStatement(genFunc, ctx, elem)
	}

	output.FIXMEDebug("generating %v %v %v", dtype, el, file)
	output.FIXMEDebug("name %v", name)

	genFunc.Emit()
}

func GenerateExternFunction(
	fp GeneratedFile,
	mod *symbols.Module,
	fn symbols.Symbol) {

	dtype := fn.Type().(symbols.FunctionDataType)
	name := MakeSymbolName(mod.Path, dtype, fn.Name()) 

	s := fmt.Sprintf("declare %v @%v(", MakeLLVMType(dtype.ReturnType()), name)
	for i,param := range dtype.Parameters() {
		if i > 0 {
			s = s + ", "
		}
		s = s + MakeLLVMType(param.DType)
	}
	s = s + ")"

	fp.Emit("%v", s)

}


func GenerateStatement(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) {

	GenerateExpression(genFunc, ctx, el)
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
	switch dtype {
	case symbols.VoidType: return "void"
	case symbols.BoolType: return "i1"
	case symbols.Int8Type: return "i8"
	case symbols.Int16Type: return "i16"
	case symbols.Int32Type: return "i32"
	case symbols.Int64Type: return "i64"
	case symbols.UInt8Type: return "i8"
	case symbols.UInt16Type: return "i16"
	case symbols.UInt32Type: return "i32"
	case symbols.UInt64Type: return "i64"
	case symbols.Real32Type: return "float"
	case symbols.Real64Type: return "double"

	//case symbols.IntegerType:
	}

	output.FatalError("Unimplemented constant type: %v", dtype)
	return "INVALID"
}

func MakeLLVMConst(val symbols.DataValue) string {

	switch val.Type() {
	case symbols.BoolType:
		if val == symbols.TrueValue {
			return "true"
		} else {
			return "false"
		}

	case symbols.Int8Type: return val.ValueAsString()
	case symbols.Int16Type: return val.ValueAsString()
	case symbols.Int32Type: return val.ValueAsString()
	case symbols.Int64Type: return val.ValueAsString()
	case symbols.UInt8Type: return val.ValueAsString()
	case symbols.UInt16Type: return val.ValueAsString()
	case symbols.UInt32Type: return val.ValueAsString()
	case symbols.UInt64Type: return val.ValueAsString()

	//FIXME better mechanism for emitting exact floating consts
	case symbols.Real32Type: return val.ValueAsString()
	case symbols.Real64Type: return val.ValueAsString()

	//case symbols.IntegerType:
	}

	output.FatalError("Unimplemented constant type: %v", val)
	return "INVALID"
}

