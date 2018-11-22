
package symbols

import (
	"fmt"
)

type CodeDataValue interface {
	EvaluateConst(op Symbol, args []DataValue) DataValue
}

type intrinsicDV struct {
	fnEvalConst func (op Symbol, args []DataValue) DataValue
}

func (self *intrinsicDV) Type() DataType {
	return IntrinsicType
}

func (self *intrinsicDV) ValueAsString() string {
	return fmt.Sprintf("@%p", self)
}

func (self *intrinsicDV) String() string {
	return DataValueString(self)
}

func (self *intrinsicDV) EvaluateConst(
	op Symbol, args []DataValue) DataValue {
	return self.fnEvalConst(op, args)
}


var IntrinsicAddInt64 = &intrinsicDV {

	fnEvalConst: func (op Symbol, args []DataValue) DataValue {

		a := args[0].(*signedDV).AsSigned64()
		b := args[1].(*signedDV).AsSigned64()

		return &signedDV{Int64Type, a+b}
	} ,

	/* fnGenerate: func (op Symbol, args []DataValue) { } */
}

var IntrinsicDivReal64 = &intrinsicDV {

	fnEvalConst: func (op Symbol, args []DataValue) DataValue {

		var a float64
		switch args[0].(type) {
			case *signedDV:
				a = float64(args[0].(*signedDV).AsSigned64())
			case *realDV:
				a = args[0].(*realDV).AsReal64()
		}

		var b float64
		switch args[1].(type) {
			case *signedDV:
				b = float64(args[1].(*signedDV).AsSigned64())
			case *realDV:
				b = args[1].(*realDV).AsReal64()
		}

		return &realDV{Real64Type, a/b}
	} ,

	/* fnGenerate: func (op Symbol, args []DataValue) { } */
}

