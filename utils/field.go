package utils

import "reflect"

func IsFieldNonEmpty(v interface{}, fieldName string) bool {
	rv := reflect.ValueOf(v)

	field := rv.FieldByName(fieldName)

	if !field.IsValid() {
		return false
	}

	return !field.IsZero()
}
