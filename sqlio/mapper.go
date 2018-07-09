package sqlio

import (
	"reflect"
	"sort"
)

// Values type alias for map[string]interface{}
type Values map[string]interface{}

// CreateMap unmarshall a struct for a CQL result
func ToMap(v *interface{}) Values {
	result := Values{}

	obj := reflect.ValueOf(*v)
	objType := obj.Type()
	numOfFields := objType.NumField()

	for i := 0; i < numOfFields; i++ {
		structField := objType.Field(i)
		if key, ok := structField.Tag.Lookup("sql"); ok {
			value := obj.FieldByName(structField.Name)
			result[key] = value.Interface()
		}
	}

	return result
}

func (values Values) SortedKeys() []string {
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
