
package symbols

import (
	"fmt"

	"parser"
)

/* FIXME clean up datavalues
	sym.Type() should always equal sym.InitialValue().Type()
		so IntrinsicType and CodeType aren't useful.

	use ValueTag to get CodeDataValue, IntrinsicDataValue etc
		add explicit DataType to intinisicDV

	FunctionChoice should be its own DataValue, with users of
	the choices selecting what they need as a conversion.  Stop
	using ctx.InitializerType for function choice selection, because
	it only works for simple assignments, while we need choice for calls,
	function reference initialization, etc.

	Create a GlobalDataRef DataValue, which includes
		the Symbol for the base object
		a byte offset
		So a GlobalDataRef has the type of the actual data being referenced
		 which may be different from the symbol being referenced
		 e.g.
		 sym{MRef(Int32)}.InitialValue() ==
		 	GlobalDataRef{
				type == MRef(Int32)
				symbol{struct whatever}
				offset{ position of Int32 struct member}
			}

*/

type ValueTag struct { string }
func (self *ValueTag) String() string {
	return self.string
}

var INTRINSIC_VALUE = &ValueTag{"INTRINSIC_VALUE"}
var CODE_VALUE = &ValueTag{"CODE_VALUE"}
var TYPE_VALUE = &ValueTag{"TYPE_VALUE"}
var BOOL_VALUE = &ValueTag{"BOOL_VALUE"}
var REAL_VALUE = &ValueTag{"REAL_VALUE"}
var SIGNED_VALUE = &ValueTag{"SIGNED_VALUE"}
var UNSIGNED_VALUE = &ValueTag{"UNSIGNED_VALUE"}
var NAMESPACE_VALUE = &ValueTag{"NAMESPACE_VALUE"}
var GLOBAL_DATA_VALUE = &ValueTag{"GLOBAL_DATA_VALUE"}
var FUNCTION_CHOICE_VALUE = &ValueTag{"FUNCTION_CHOICE_VALUE"}


type DataValue interface {
	Tag() *ValueTag
	Type() DataType
	ValueAsString() string
	String() string
}

func DataValueString(dv DataValue) string {
	return dv.Type().String() + "{" + dv.ValueAsString() + "}"
}

type RealDataValue interface {
	DataValue
	AsReal64() float64
}

type realDV struct {
	dtype DataType
	value float64
}

func (self *realDV) Tag() *ValueTag {
	return REAL_VALUE
}

func (self *realDV) Type() DataType {
	return self.dtype
}

func (self *realDV) ValueAsString() string {
	return fmt.Sprintf("%e", self.value)
}

func (self *realDV) String() string {
	return DataValueString(self)
}

func (self *realDV) AsReal64() float64 {
	return self.value
}

type SignedDataValue interface {
	DataValue
	AsSigned64() int64
}

type signedDV struct {
	dtype DataType
	value int64
}

func (self *signedDV) Tag() *ValueTag {
	return SIGNED_VALUE
}

func (self *signedDV) Type() DataType {
	return self.dtype
}

func (self *signedDV) ValueAsString() string {
	return fmt.Sprintf("%v", self.value)
}

func (self *signedDV) String() string {
	return DataValueString(self)
}

func (self *signedDV) AsSigned64() int64 {
	return self.value
}

type UnsignedDataValue interface {
	DataValue
	AsUnsigned64() uint64
}

type unsignedDV struct {
	dtype DataType
	value uint64
}

func (self *unsignedDV) Tag() *ValueTag {
	return UNSIGNED_VALUE
}

func (self *unsignedDV) Type() DataType {
	return self.dtype
}

func (self *unsignedDV) ValueAsString() string {
	return fmt.Sprintf("%v", self.value)
}

func (self *unsignedDV) String() string {
	return DataValueString(self)
}

func (self *unsignedDV) AsUnsigned64() uint64 {
	return self.value
}


type BoolDataValue interface {
	DataValue
	AsBool() bool
}

type boolDV struct {
	value bool
}

func (self *boolDV) Tag() *ValueTag {
	return BOOL_VALUE
}

func (self *boolDV) Type() DataType {
	return BoolType
}

func (self *boolDV) ValueAsString() string {
	if self.value {
		return "TRUE"
	} else {
		return "FALSE"
	}
}

func (self *boolDV) String() string {
	return DataValueString(self)
}

func (self *boolDV) AsBool() bool {
	return self.value
}


type NamespaceDataValue interface {
	DataValue
	AsSymbolTable() *symbolTable
}

type namespaceDV struct {
	value *symbolTable
}

func (self *namespaceDV) Tag() *ValueTag {
	return NAMESPACE_VALUE
}

func (self *namespaceDV) Type() DataType {
	return NamespaceType
}

func (self *namespaceDV) ValueAsString() string {
	return self.value.Name
}

func (self *namespaceDV) String() string {
	return DataValueString(self)
}

func (self *namespaceDV) AsSymbolTable() *symbolTable {
	return self.value
}

type TypeDataValue interface {
	DataValue
	AsDataType() DataType
}

type typeDV struct {
	dtype DataType
	value DataType
}

func (self *typeDV) Tag() *ValueTag {
	return TYPE_VALUE
}

func (self *typeDV) Type() DataType {
	return self.dtype
}

func (self *typeDV) ValueAsString() string {
	return self.value.String()
}

func (self *typeDV) String() string {
	return DataValueString(self)
}

func (self *typeDV) AsDataType() DataType {
	return self.value
}

type CodeDataValue interface {
	DataValue
	AsParseElement() parser.ParseElement
	AsSourceFile() *SourceFile
}

type codeDV struct {
	dtype DataType
	element parser.ParseElement
	file *SourceFile
}

func (self *codeDV) Tag() *ValueTag {
	return CODE_VALUE
}

func (self *codeDV) Type() DataType {
	return self.dtype
}

func (self *codeDV) ValueAsString() string {
	return self.element.ElementType().String()
}

func (self *codeDV) String() string {
	return DataValueString(self)
}

func (self *codeDV) AsParseElement() parser.ParseElement {
	return self.element
}

func (self *codeDV) AsSourceFile() *SourceFile {
	return self.file
}

type GlobalDataRefValue interface {
	DataValue
	AsSymbol() Symbol
	AsOffset() int64
}
type globalDV struct {
 	//dtype will be Ref(symbol.Type) or MRef(symbol.Type())
	// except for functions where dtype = symbol.Type()
	dtype DataType
	symbol Symbol
	offset int64
}

func (self *globalDV) Tag() *ValueTag { return GLOBAL_DATA_VALUE }
func (self *globalDV) Type() DataType { return self.dtype }
func (self *globalDV) AsSymbol() Symbol { return self.symbol }
func (self *globalDV) AsOffset() int64 { return self.offset }

func (self *globalDV) ValueAsString() string {
	return fmt.Sprintf("global_ref(%v)", self.symbol)
}

func (self *globalDV) String() string {
	return DataValueString(self)
}

type FunctionChoiceValue interface {
	DataValue
	AsSymbol() FunctionChoiceSymbol
}

type functionChoiceDV struct {
	symbol FunctionChoiceSymbol
}

func (self *functionChoiceDV) Tag() *ValueTag { return FUNCTION_CHOICE_VALUE }
func (self *functionChoiceDV) Type() DataType { return FunctionChoiceType }
func (self *functionChoiceDV) AsSymbol() FunctionChoiceSymbol {
	return self.symbol
}

func (self *functionChoiceDV) ValueAsString() string {
	return fmt.Sprintf("choice(%v)", self.symbol)
}

func (self *functionChoiceDV) String() string {
	return DataValueString(self)
}

