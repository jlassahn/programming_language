
package symbols



var predefinedSymbols *symbolTable;

func PredefinedSymbols() *symbolTable {
	if predefinedSymbols == nil {
		predefinedSymbols = buildPredefinedSymbols()
	}
	return predefinedSymbols
}


//FIXME reorganize and correct
func buildPredefinedSymbols() *symbolTable {

	syms := NewSymbolTable("PREDEFINED", nil)

	syms.Symbols["__system"] = buildInternalSymbols()

	syms.AddConst("TRUE", BoolType, TrueValue)
	syms.AddConst("FALSE", BoolType, FalseValue)

	syms.AddConst("CType", CTypeType, &typeDV{CTypeType, CTypeType})
	syms.AddConst("Bool", CTypeType, &typeDV{CTypeType, BoolType})
	syms.AddConst("Int32", CTypeType, &typeDV{CTypeType, Int32Type})

	addBinaryIntrinsic(syms, "+", "add_Int64", Int64Type)
	addBinaryIntrinsic(syms, "+", "add_Int32", Int32Type)
	addBinaryIntrinsic(syms, "+", "add_Int16", Int16Type)
	addBinaryIntrinsic(syms, "+", "add_Int8", Int8Type)
	addBinaryIntrinsic(syms, "+", "add_UInt64", UInt64Type)
	addBinaryIntrinsic(syms, "+", "add_UInt32", UInt32Type)
	addBinaryIntrinsic(syms, "+", "add_UInt16", UInt16Type)
	addBinaryIntrinsic(syms, "+", "add_UInt8", UInt8Type)
	addBinaryIntrinsic(syms, "+", "add_Real32", Real32Type)
	addBinaryIntrinsic(syms, "+", "add_Real64", Real64Type)

	syms.AddOperator("/", Real64Type, []FunctionParameter {
		{"a", Real64Type, false},
		{"b", Real64Type, true},
	},
	true, &intrinsicDV{"div_Real64"})

	//FIXME should  be Int64, >Real64
	//FIXME should have other IntXXX versions
	syms.AddOperator("/", Real64Type, []FunctionParameter {
		{"a", Int64Type, false},
		{"b", Int64Type, true},
	},
	true, &intrinsicDV{"div_Real64"})

	return syms
}

func buildInternalSymbols() *baseSymbol {

	name := "PREDEFINED:__system"
	newTable := NewSymbolTable(name, nil)

	val := &namespaceDV {
		value: newTable,
	}

	return &baseSymbol {
		name: name,
		dtype: NamespaceType,
		initialValue: val,
		isConst: true,
	}
}

func addBinaryIntrinsic(syms *symbolTable, name string, op string,
	dtype DataType) {

	syms.AddOperator(name, dtype, []FunctionParameter {
		{"a", dtype, false},
		{"b", dtype, true},
	},
	true, &intrinsicDV{op} )
}

