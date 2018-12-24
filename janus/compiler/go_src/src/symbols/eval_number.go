
package symbols

import (
	"parser"
)

const CAT_REAL = 1
const CAT_SIGNED = 2
const CAT_UNSIGNED = 3
const CAT_BIGNUM = 4

type typeInfo struct {
	dtype DataType
	cat int
}

var typeInfoFromTag = map[string] *typeInfo {

	"i" : {IntegerType, CAT_BIGNUM},
	"u" : {UInt64Type, CAT_UNSIGNED},
	"s" : {Int64Type, CAT_SIGNED},

	"r" : {Real64Type, CAT_REAL},

	"r64" : {Real64Type, CAT_REAL},
	"r32" : {Real32Type, CAT_REAL},

	"s64" : {Int64Type, CAT_SIGNED},
	"s32" : {Int32Type, CAT_SIGNED},
	"s16" : {Int16Type, CAT_SIGNED},
	"s8" : {Int8Type, CAT_SIGNED},

	"u64" : {UInt64Type, CAT_UNSIGNED},
	"u32" : {UInt32Type, CAT_UNSIGNED},
	"u16" : {UInt16Type, CAT_UNSIGNED},
	"u8" : {UInt8Type, CAT_UNSIGNED},

}


func evalNumber(el parser.ParseElement, ctx *EvalContext) DataValue {

	pos := el.FilePos()
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
			parser.Error(pos, "invalid character in numeric constant")
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
				parser.Error(pos, "invalid character in numeric constant")
				continue
			}

			mult = mult/float64(base)
			frac = frac + mult*float64(digit)
		}
		if mult == 1.0 {
			parser.Error(pos, "missing digit after decimal point")
		}
	}

	tag := txt[i:]

	if tag == "" {
		if is_float {
			return &realDV{Real64Type, float64(iv) + frac}
		} else {
			return &signedDV{Int64Type, int64(iv)}
		}
	}

	typeInfo := typeInfoFromTag[tag]
	if typeInfo == nil {
		parser.Error(pos, "invalid type specifier for number constant")
		return nil
	}

	switch typeInfo.cat {
		case CAT_REAL:
			return &realDV{typeInfo.dtype, float64(iv) + frac}

		case CAT_SIGNED:
			if is_float {
				parser.Error(pos, "fractional part in integer constant")
				return nil
			}

			return &signedDV{typeInfo.dtype, int64(iv)}

		case CAT_UNSIGNED:
			if is_float {
				parser.Error(pos, "fractional part in integer constant")
				return nil
			}
			return &unsignedDV{typeInfo.dtype, iv}

		case CAT_BIGNUM:
			parser.Error(pos, "FIXME large integer values not implemented")
			return nil
	}

	//can't happen
	return nil
}

