
package symbols

type DataType struct {
}

type DataValue struct {
}

type Symbol struct {
	data_type *DataType
	initial_value *DataValue
}

func ResolveGlobals(file_set *FileSet) {
	//FIXME implement
}

