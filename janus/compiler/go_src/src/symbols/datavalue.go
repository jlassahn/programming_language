
package symbols

import (
	"fmt"
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

