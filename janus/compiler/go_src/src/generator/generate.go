
package generator

import (
	"fmt"
	"bytes"
	"encoding/binary"

	"output"
	"parser"
	"symbols"
)


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
		genFile.Emit("@janus_main = alias void(), void()* @%v", mainName)
	}

}

func GenerateHeader(fp GeneratedFile) {
	fp.EmitComment("")
	fp.EmitComment("global declarations")
	fp.EmitComment("")

	//FIXME organize better
	fp.Emit("declare double @llvm.sqrt.f64(double)")
	fp.Emit("declare void @clib_print_Real64(double)")
	fp.Emit("declare void @clib_print_Int64(i64)")

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

	case symbols.Real32Type: return MakeLLVMReal(val)
	case symbols.Real64Type: return MakeLLVMReal(val)

	//case symbols.IntegerType:
	}

	output.FatalError("Unimplemented constant type: %v", val)
	return "INVALID"
}

func MakeLLVMReal(dval symbols.DataValue) string {
	x := dval.(symbols.RealDataValue).AsReal64()
	if dval.Type() == symbols.Real32Type {
		x = float64(float32(x))
	}
	return MakeLLVMDouble(x)
}

func MakeLLVMDouble(x float64) string {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, x)
	var n uint64
	binary.Read(buf, binary.LittleEndian, &n)
	return fmt.Sprintf("0x%x", n)
}

