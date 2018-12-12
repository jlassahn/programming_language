
package symbols

import (
	"sort"
	"reflect"
)


func SortedKeys(stringMap interface{}) []string {

	val := reflect.ValueOf(stringMap)
	keyVals := val.MapKeys()

	keys := make([]string, len(keyVals))
	for i,x := range keyVals {
		keys[i] = x.String()
	}

	sort.Strings(keys)
	return keys
}

