
package symbols

import (
	"output"
	"lexer"
	"parser"
	"fmt"
)


type Evaluator interface {
	EvaluateConstExpression(el parser.ParseElement, ctx *SymbolTable) *DataValue
}

var evaluators = map[lexer.Tag] Evaluator {
	*lexer.NUMBER: &NumberEval {} }


type NumberEval struct {}
func (*NumberEval) EvaluateConstExpression(
		el parser.ParseElement, ctx *SymbolTable) *DataValue {

	line, col := el.Position()
	txt := el.TokenString()

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
		fmt.Printf("digit=%v iv=%v\n", digit, iv)
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

	switch txt[i:] {
	case "":
	case "i":
	case "u":
	case "u8":
	case "u16":
	case "u32":
	case "u64":
	case "s":
	case "s8":
	case "s16":
	case "s32":
	case "s64":
	case "r":
	case "r32":
	case "r64":

	default:
		output.Error(line, col+i, "invalid character in numeric constant")
	}

	fmt.Printf(" iv = %v frac = %v is_float = %v tail = %v\n",
		iv, frac, is_float, txt[i:])

	return nil
}

func DotListAsStrings(el parser.ParseElement) []string {
	var ret []string
	for _, x := range(el.Children()) {
		ret = append(ret, x.TokenString())
	}
	return ret
}

func EvaluateConstExpression(
	el parser.ParseElement, ctx *SymbolTable) *DataValue {

	eval := evaluators[*el.ElementType()]
	if eval == nil {
		//FIXME implement
		parser.EmitParseTree(el)
		EmitSymbolTable(ctx)
		return nil
	} else {
		return eval.EvaluateConstExpression(el, ctx)
	}
}

