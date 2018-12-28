
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
	syms.AddConst("Int8", CTypeType, &typeDV{CTypeType, Int8Type})
	syms.AddConst("Int16", CTypeType, &typeDV{CTypeType, Int16Type})
	syms.AddConst("Int32", CTypeType, &typeDV{CTypeType, Int32Type})
	syms.AddConst("Int64", CTypeType, &typeDV{CTypeType, Int64Type})
	syms.AddConst("UInt8", CTypeType, &typeDV{CTypeType, UInt8Type})
	syms.AddConst("UInt16", CTypeType, &typeDV{CTypeType, UInt16Type})
	syms.AddConst("UInt32", CTypeType, &typeDV{CTypeType, UInt32Type})
	syms.AddConst("UInt64", CTypeType, &typeDV{CTypeType, UInt64Type})
	syms.AddConst("Real32", CTypeType, &typeDV{CTypeType, Real32Type})
	syms.AddConst("Real64", CTypeType, &typeDV{CTypeType, Real64Type})

	addUnaryIntrinsic(syms, "-", "negate_Int8", Int8Type)
	addUnaryIntrinsic(syms, "-", "negate_Int16", Int16Type)
	addUnaryIntrinsic(syms, "-", "negate_Int32", Int32Type)
	addUnaryIntrinsic(syms, "-", "negate_Int64", Int64Type)
	addUnaryIntrinsic(syms, "-", "negate_Real32", Real32Type)
	addUnaryIntrinsic(syms, "-", "negate_Real64", Real64Type)

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

	addBinaryIntrinsic(syms, "-", "sub_Int64", Int64Type)
	addBinaryIntrinsic(syms, "-", "sub_Int32", Int32Type)
	addBinaryIntrinsic(syms, "-", "sub_Int16", Int16Type)
	addBinaryIntrinsic(syms, "-", "sub_Int8", Int8Type)
	addBinaryIntrinsic(syms, "-", "sub_UInt64", UInt64Type)
	addBinaryIntrinsic(syms, "-", "sub_UInt32", UInt32Type)
	addBinaryIntrinsic(syms, "-", "sub_UInt16", UInt16Type)
	addBinaryIntrinsic(syms, "-", "sub_UInt8", UInt8Type)
	addBinaryIntrinsic(syms, "-", "sub_Real32", Real32Type)
	addBinaryIntrinsic(syms, "+", "sub_Real64", Real64Type)

	addBinaryIntrinsic(syms, "/", "div_Real64", Real64Type)
	addBinaryIntrinsic(syms, "/", "div_Real32", Real32Type)
	syms.AddOperator("/", Real64Type, []FunctionParameter {
		{"a", Int64Type, true},
		{"b", Real64Type, true},
	},
	true, &intrinsicDV{"div_IntReal"})

	return syms
}

func buildInternalSymbols() *baseSymbol {

	name := "PREDEFINED:__system"
	newTable := NewSymbolTable(name, nil)

	newTable.AddFunction("sqrt", Real64Type, []FunctionParameter {
		{"a", Real64Type, false}, //FIXME should be true
	} ,
	true, &intrinsicDV{"sqrt_Real64"})

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

func addUnaryIntrinsic(syms *symbolTable, name string, op string,
	dtype DataType) {

	syms.AddOperator(name, dtype, []FunctionParameter {
		{"a", dtype, false},
	},
	true, &intrinsicDV{op} )
}

