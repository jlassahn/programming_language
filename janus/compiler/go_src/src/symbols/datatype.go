
package symbols

import (
)

type Tag struct { string }

type DataType interface {
	String() string
	Base() *Tag
	SubTypes() []DataValue
}

type SimpleDataType struct {
	Tag
}

func (self *SimpleDataType) Base() *Tag {
	return &self.Tag
}

func (self *SimpleDataType) SubTypes() []DataValue {
	return nil //FIXME empty array instead?
}

func (self *SimpleDataType) String() string {
	return self.string
}

//FIXME make this a member??

func TypeString(dt DataType) string {
	ret := dt.Base().string
	if dt.SubTypes != nil {
		ret += "("
		for i, st := range dt.SubTypes() {
			if i > 0 {
				ret += ", "
			}
			ret += st.ValueAsString()
		}
		ret += ")"
	}

	return ret
}

