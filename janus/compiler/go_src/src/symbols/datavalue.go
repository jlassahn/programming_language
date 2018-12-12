
package symbols

import (
	"fmt"

	"parser"
)

type DataValue interface {
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

func (self *realDV) Type() DataType {
	return self.dtype
}

func (self *realDV) ValueAsString() string {
	return fmt.Sprintf("%v", self.value)
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

