
package generator

import (
	"fmt"

	"output"
	"symbols"
)


func MakeIntrinsicOp(ret Result, opName string, args []Result) string {

	op, ok := LLVMOperator[opName]
	if ok {
		return  fmt.Sprintf("\t%v = %v %v %v, %v",
			ret.LLVMVal(),
			op,
			args[0].LLVMType(),
			args[0].LLVMVal(),
			args[1].LLVMVal())
	}

	//FIXME this could be done by using the normal function call path
	//   and providing a Result with name = LLVMFunction[opName]
	//   which could be injected into the intrinsic by SetGenVal
	op, ok = LLVMFunction[opName]
	if ok {

		var s string
		if ret.Type() == symbols.VoidType {
			s = fmt.Sprintf("\tcall %v %v(",
				ret.LLVMType(),
				op)
		} else {
			s = fmt.Sprintf("\t%v = call %v %v(",
				ret.LLVMVal(),
				ret.LLVMType(),
				op)
		}

		for i,arg := range args {
			if i > 0 {
				s = s + ", "
			}
			s = s + fmt.Sprintf("%v %v", arg.LLVMType(), arg.LLVMVal())
		}
		s = s + ")"

		return s
	}

	output.FatalError("Unimplemented intrinsic %v", opName)

	return ""
}

func MakeLLVMConvert(from Result, ret Result) string {

	to := ret.Type()

	opString, ok := baseTypeConvert[tagPair{from.Type().Base(), to.Base()}]
	if ok {
		return fmt.Sprintf("\t%v = %s %v %v to %v",
			ret.LLVMVal(),
			opString,
			from.LLVMType(),
			from.LLVMVal(),
			ret.LLVMType())
	}

	return ""
}

var LLVMOperator = map[string]string {
	//FIXME add a bunch of stuff here
	"add_Int8": "add",
	"add_Int16": "add",
	"add_Int32": "add",
	"add_Int64": "add",
	"add_UInt8": "add",
	"add_UInt16": "add",
	"add_UInt32": "add",
	"add_UInt64": "add",
	"add_Real32": "fadd",
	"add_Real64": "fadd",

	"sub_Int8": "sub",
	"sub_Int16": "sub",
	"sub_Int32": "sub",
	"sub_Int64": "sub",
	"sub_UInt8": "sub",
	"sub_UInt16": "sub",
	"sub_UInt32": "sub",
	"sub_UInt64": "sub",
	"sub_Real32": "fsub",
	"sub_Real64": "fsub",

	"mul_Int8": "mul",
	"mul_Int16": "mul",
	"mul_Int32": "mul",
	"mul_Int64": "mul",
	"mul_UInt8": "mul",
	"mul_UInt16": "mul",
	"mul_UInt32": "mul",
	"mul_UInt64": "mul",
	"mul_Real32": "fmul",
	"mul_Real64": "fmul",

	"cmp_eq_Int8": "icmp eq",
	"cmp_eq_Int16": "icmp eq",
	"cmp_eq_Int32": "icmp eq",
	"cmp_eq_Int64": "icmp eq",
	"cmp_ne_Int8": "icmp ne",
	"cmp_ne_Int16": "icmp ne",
	"cmp_ne_Int32": "icmp ne",
	"cmp_ne_Int64": "icmp ne",
	"cmp_lt_Int8": "icmp slt",
	"cmp_lt_Int16": "icmp slt",
	"cmp_lt_Int32": "icmp slt",
	"cmp_lt_Int64": "icmp slt",
	"cmp_le_Int8": "icmp sle",
	"cmp_le_Int16": "icmp sle",
	"cmp_le_Int32": "icmp sle",
	"cmp_le_Int64": "icmp sle",
	"cmp_ge_Int8": "icmp sge",
	"cmp_ge_Int16": "icmp sge",
	"cmp_ge_Int32": "icmp sge",
	"cmp_ge_Int64": "icmp sge",
	"cmp_gt_Int8": "icmp sgt",
	"cmp_gt_Int16": "icmp sgt",
	"cmp_gt_Int32": "icmp sgt",
	"cmp_gt_Int64": "icmp sgt",

	"cmp_eq_Real64": "fcmp oeq",
	"cmp_ne_Real64": "fcmp une",
	"cmp_lt_Real64": "fcmp olt",
	"cmp_le_Real64": "fcmp ole",
	"cmp_ge_Real64": "fcmp oge",
	"cmp_gt_Real64": "fcmp ogt",
	"cmp_eq_Real32": "fcmp oeq",
	"cmp_ne_Real32": "fcmp une",
	"cmp_lt_Real32": "fcmp olt",
	"cmp_le_Real32": "fcmp ole",
	"cmp_ge_Real32": "fcmp oge",
	"cmp_gt_Real32": "fcmp ogt",

}

var LLVMFunction = map[string]string {
	"sqrt_Real64": "@llvm.sqrt.f64",
	"print_Real64": "@clib_print_Real64",
	"print_Int64": "@clib_print_Int64",
}

//FIXME reorganize
type tagPair struct {
	from *symbols.Tag
	to *symbols.Tag
}

var baseTypeConvert = map[tagPair] string  {
	{symbols.INT8_TYPE, symbols.INT16_TYPE}: "sext",
	{symbols.INT8_TYPE, symbols.INT32_TYPE}: "sext",
	{symbols.INT8_TYPE, symbols.INT64_TYPE}: "sext",
	{symbols.INT16_TYPE, symbols.INT32_TYPE}: "sext",
	{symbols.INT16_TYPE, symbols.INT64_TYPE}: "sext",
}

