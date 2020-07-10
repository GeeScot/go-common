package env

import (
	"os"
	"reflect"
)

// Read environment variables in to struct, fields must have 'env' tag.
func Read(envStruct interface{}) {
	reflectValue := reflect.ValueOf(envStruct)
	reflectElem := reflectValue.Elem()
	reflectType := reflectElem.Type()

	numOfFields := reflectType.NumField()
	for i := 0; i < numOfFields; i++ {
		structField := reflectType.Field(i)
		if key, ok := structField.Tag.Lookup("env"); ok {
			fieldValue := reflectElem.FieldByName(structField.Name)
			if fieldValue.CanAddr() && fieldValue.CanSet() {
				value := os.Getenv(key)
				fieldValue.SetString(value)
			}
		}
	}
}

// Optional attempt to get an optional environment variable or choose default value
func Optional(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
