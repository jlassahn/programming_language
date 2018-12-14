
package symbols

import (
	"output"
)

func EvaluateIntrinsic(opName string, args []DataValue) DataValue {

	fn := intrinsics[opName]
	if fn == nil {
		output.FatalError("Unimplemented intrinsic: %v", opName)
		return nil
	}

	return fn(opName, args)
}

var intrinsics = map[string] func(name string, args []DataValue) DataValue {
	"add_Int64": intrinsicAddInt64,
	"div_Real64": intrinsicDivReal64,
}

type IntrinsicDataValue interface {
	ValueAsString() string
}

type intrinsicDV struct {
	name string
}

func (self *intrinsicDV) Type() DataType {
	return IntrinsicType
}

func (self *intrinsicDV) ValueAsString() string {
	return self.name
}

func (self *intrinsicDV) String() string {
	return DataValueString(self)
}

func intrinsicAddInt64(op string, args []DataValue) DataValue {

	a := args[0].(*signedDV).AsSigned64()
	b := args[1].(*signedDV).AsSigned64()

	return &signedDV{Int64Type, a+b}
}

//FIXME should be real args, with conversions
func intrinsicDivReal64(op string, args []DataValue) DataValue {

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
}

