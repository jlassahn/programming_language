
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

	//FIXME where should this live
	marray := &paramDT{
		tag: MARRAY_TYPE,
		params: []DTypeParameter{
			{0, &typevarDT{1, false}},
			{0, &typevarDT{2, true}},
		},
		members: nil, //FIXME add conversions, etc
	}
	syms.AddConst("MArray", MetaTypeType, &typeDV{MetaTypeType, marray})

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
	addBinaryIntrinsic(syms, "-", "sub_Real64", Real64Type)

	addBinaryIntrinsic(syms, "*", "mul_Int64", Int64Type)
	addBinaryIntrinsic(syms, "*", "mul_Int32", Int32Type)
	addBinaryIntrinsic(syms, "*", "mul_Int16", Int16Type)
	addBinaryIntrinsic(syms, "*", "mul_Int8", Int8Type)
	addBinaryIntrinsic(syms, "*", "mul_UInt64", UInt64Type)
	addBinaryIntrinsic(syms, "*", "mul_UInt32", UInt32Type)
	addBinaryIntrinsic(syms, "*", "mul_UInt16", UInt16Type)
	addBinaryIntrinsic(syms, "*", "mul_UInt8", UInt8Type)
	addBinaryIntrinsic(syms, "*", "mul_Real32", Real32Type)
	addBinaryIntrinsic(syms, "*", "mul_Real64", Real64Type)

	addBinaryIntrinsic(syms, "/", "div_Real64", Real64Type)
	addBinaryIntrinsic(syms, "/", "div_Real32", Real32Type)

	//FIXME figure out how to handle division
	fnType := NewFunction(Real64Type)
	fnType.AddParam("a", Int64Type, true)
	fnType.AddParam("b", Real64Type, true)
	syms.AddOperator("/", fnType, true, &intrinsicDV{fnType, "div_IntReal"})

	addCompareIntrinsic(syms, "==", "cmp_eq_Int8", Int8Type)
	addCompareIntrinsic(syms, "!=", "cmp_ne_Int8", Int8Type)
	addCompareIntrinsic(syms, "<",  "cmp_lt_Int8", Int8Type)
	addCompareIntrinsic(syms, "<=", "cmp_le_Int8", Int8Type)
	addCompareIntrinsic(syms, ">=", "cmp_ge_Int8", Int8Type)
	addCompareIntrinsic(syms, ">",  "cmp_gt_Int8", Int8Type)
	addCompareIntrinsic(syms, "==", "cmp_eq_Int16", Int16Type)
	addCompareIntrinsic(syms, "!=", "cmp_ne_Int16", Int16Type)
	addCompareIntrinsic(syms, "<",  "cmp_lt_Int16", Int16Type)
	addCompareIntrinsic(syms, "<=", "cmp_le_Int16", Int16Type)
	addCompareIntrinsic(syms, ">=", "cmp_ge_Int16", Int16Type)
	addCompareIntrinsic(syms, ">",  "cmp_gt_Int16", Int16Type)
	addCompareIntrinsic(syms, "==", "cmp_eq_Int32", Int32Type)
	addCompareIntrinsic(syms, "!=", "cmp_ne_Int32", Int32Type)
	addCompareIntrinsic(syms, "<",  "cmp_lt_Int32", Int32Type)
	addCompareIntrinsic(syms, "<=", "cmp_le_Int32", Int32Type)
	addCompareIntrinsic(syms, ">=", "cmp_ge_Int32", Int32Type)
	addCompareIntrinsic(syms, ">",  "cmp_gt_Int32", Int32Type)
	addCompareIntrinsic(syms, "==", "cmp_eq_Int64", Int64Type)
	addCompareIntrinsic(syms, "!=", "cmp_ne_Int64", Int64Type)
	addCompareIntrinsic(syms, "<",  "cmp_lt_Int64", Int64Type)
	addCompareIntrinsic(syms, "<=", "cmp_le_Int64", Int64Type)
	addCompareIntrinsic(syms, ">=", "cmp_ge_Int64", Int64Type)
	addCompareIntrinsic(syms, ">",  "cmp_gt_Int64", Int64Type)
	addCompareIntrinsic(syms, "==", "cmp_eq_Real32", Real32Type)
	addCompareIntrinsic(syms, "!=", "cmp_ne_Real32", Real32Type)
	addCompareIntrinsic(syms, "<",  "cmp_lt_Real32", Real32Type)
	addCompareIntrinsic(syms, "<=", "cmp_le_Real32", Real32Type)
	addCompareIntrinsic(syms, ">=", "cmp_ge_Real32", Real32Type)
	addCompareIntrinsic(syms, ">",  "cmp_gt_Real32", Real32Type)
	addCompareIntrinsic(syms, "==", "cmp_eq_Real64", Real64Type)
	addCompareIntrinsic(syms, "!=", "cmp_ne_Real64", Real64Type)
	addCompareIntrinsic(syms, "<",  "cmp_lt_Real64", Real64Type)
	addCompareIntrinsic(syms, "<=", "cmp_le_Real64", Real64Type)
	addCompareIntrinsic(syms, ">=", "cmp_ge_Real64", Real64Type)
	addCompareIntrinsic(syms, ">",  "cmp_gt_Real64", Real64Type)

	return syms
}

func buildInternalSymbols() *baseSymbol {

	name := "PREDEFINED:__system"
	newTable := NewSymbolTable(name, nil)

	fnType := NewFunction(Real64Type)
	fnType.AddParam("a", Real64Type, false) //FIXME should be true
	newTable.AddFunction("sqrt", fnType, true, &intrinsicDV{fnType, "sqrt_Real64"})

	//FIXME ---- fake print statements ----
	fnType = NewFunction(VoidType)
	fnType.AddParam("a", Real64Type, false)
	newTable.AddFunction("print", fnType, false,
		&intrinsicDV{fnType, "print_Real64"})

	fnType = NewFunction(VoidType)
	fnType.AddParam("a", Int64Type, false)
	newTable.AddFunction("print", fnType, false,
		&intrinsicDV{fnType, "print_Int64"})
	// ------------------------------------

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

	fnType := NewFunction(dtype)
	fnType.AddParam("a", dtype, false)
	fnType.AddParam("b", dtype, true)

	syms.AddOperator(name, fnType, true, &intrinsicDV{fnType, op})
}

func addUnaryIntrinsic(syms *symbolTable, name string, op string,
	dtype DataType) {

	fnType := NewFunction(dtype)
	fnType.AddParam("a", dtype, false)

	syms.AddOperator(name, fnType, true, &intrinsicDV{fnType, op})
}

func addCompareIntrinsic(syms *symbolTable, name string, op string,
	dtype DataType) {

	fnType := NewFunction(BoolType)
	fnType.AddParam("a", dtype, false)
	fnType.AddParam("b", dtype, true)

	syms.AddOperator(name, fnType, true, &intrinsicDV{fnType, op})
}

