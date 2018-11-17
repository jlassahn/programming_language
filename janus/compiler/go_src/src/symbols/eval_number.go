
package symbols

import (
	"output"
	"parser"
)

type PTKey struct {
	tag string
	IsFloat bool
}

var preferredTypes = map[PTKey] []DataType {

	PTKey{"", false}:
		{Int64Type, Int32Type, Int16Type, Int8Type, IntegerType,
		Real64Type, Real32Type},

	PTKey{"", true}: {Real64Type, Real32Type},

	PTKey{"i", false} : {IntegerType},
	PTKey{"u", false} : {UInt64Type, UInt32Type, UInt16Type, UInt8Type},
	PTKey{"s", false} : {Int64Type, Int32Type, Int16Type, Int8Type},

	PTKey{"r", false} : {Real64Type, Real32Type},
	PTKey{"r", true} : {Real64Type, Real32Type},

	PTKey{"r64", false} : {Real64Type},
	PTKey{"r64", true} : {Real64Type},
	PTKey{"r32", false} : {Real32Type},
	PTKey{"r32", true} : {Real32Type},

	PTKey{"s64", false} : {Int64Type},
	PTKey{"s32", false} : {Int32Type},
	PTKey{"s16", false} : {Int16Type},
	PTKey{"s8", false} : {Int8Type},

	PTKey{"u64", false} : {UInt64Type},
	PTKey{"u32", false} : {UInt32Type},
	PTKey{"u16", false} : {UInt16Type},
	PTKey{"u8", false} : {UInt8Type},

}

const CAT_REAL = 1
const CAT_SIGNED = 2
const CAT_UNSIGNED = 3
const CAT_BIGNUM = 4

var categoryForType = map[DataType] int {
	IntegerType: CAT_BIGNUM,
	Int8Type: CAT_SIGNED,
	Int16Type: CAT_SIGNED,
	Int32Type: CAT_SIGNED,
	Int64Type: CAT_SIGNED,
	UInt8Type: CAT_UNSIGNED,
	UInt16Type: CAT_UNSIGNED,
	UInt32Type: CAT_UNSIGNED,
	UInt64Type: CAT_UNSIGNED,
	Real32Type: CAT_REAL,
	Real64Type: CAT_REAL,
}

type NumberEval struct {}
func (*NumberEval) EvaluateConstExpression(
		el parser.ParseElement, ctx *EvalContext) DataValue {

	line, col := el.Position()
	txt := el.TokenString()

	//FIXME handle bignum values

	var iv uint64 = 0
	i := 0
	base := 10
	if len(txt) > 2 && txt[0] == '0' {
		if txt[1] == 'x' { base = 16; i = 2; }
		if txt[1] == 'o' { base = 8; i = 2; }
		if txt[1] == 'b' { base = 2; i = 2; }
		if txt[1] == 'd' { base = 10; i = 2; }
	}

	for ; i < len(txt); i++ {

		c := txt[i]

		if c == '_' {
			continue
		}

		digit := 0
		if c >= '0' && c <= '9' {
			digit = int(c) - '0'
		} else if int(c) >= 'a' && int(c) <= 'f' {
			digit = int(c) - 'a' + 10
		} else if c >= 'A' && c <= 'F' {
			digit = int(c) - 'A' + 10
		} else {
			break
		}

		if digit >= base {
			output.Error(line, col+i, "invalid character in numeric constant")
			continue
		}

		iv = iv*uint64(base) + uint64(digit)
	}

	frac := 0.0
	is_float := false
	if len(txt) > i && txt[i] == '.' {
		i ++
		is_float = true
		mult := 1.0

		for ; i<len(txt); i++ {

			c := txt[i]

			if c == '_' {
				continue
			}

			digit := 0
			if c >= '0' && c <= '9' {
				digit = int(c) - '0'
			} else if c >= 'a' && c <= 'f' {
				digit = int(c) - 'a' + 10
			} else if c >= 'A' && c <= 'F' {
				digit = int(c) - 'A' + 10
			} else {
				break
			}

			if digit >= base {
				output.Error(line, col+i, "invalid character in numeric constant")
				continue
			}

			mult = mult/float64(base)
			frac = frac + mult*float64(digit)
		}
		if mult == 1.0 {
			output.Error(line, col+i, "missing digit after decimal point")
		}
	}

	tag := txt[i:]

	types := preferredTypes[PTKey{tag, is_float}]

	if types == nil {
		output.Error(line, col, "invalid type specifier for number constant")
		return nil
	}

	bestType := types[0]
	for _, t := range types {
		if t == ctx.PreferredType {
			bestType = t
			break
		}
	}

	//FIXME bounds check values

	switch categoryForType[bestType] {
		case CAT_REAL:
			return &realDV{bestType, float64(iv) + frac}

		case CAT_SIGNED:
			return &signedDV{bestType, int64(iv)}

		case CAT_UNSIGNED:
			return &unsignedDV{bestType, iv}

		case CAT_BIGNUM:
			output.Error(line, col, "FIXME large integer values not implemented")
			return nil
	}

	//can't happen
	return nil
}

