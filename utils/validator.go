package utils

import (
	"fmt"
	"reflect"
	"strings"
)

type validate struct {
	fieldvalue interface{}
	validator  string
}

func Validate(reqs any) (bool, map[string]string) {
	var validatorMessage = map[string]string{}

	t := reflect.TypeOf(reqs)
	v := reflect.ValueOf(reqs)

	// fmt.Println(t)
	// 1. Inspect fields

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		fmt.Println("Error: Parameter yang di kirim bukan struct")
		return false, validatorMessage
	}

	var validatorSlice = map[string]validate{}

	for i := 0; i < t.NumField(); i++ {

		fieldInfo := t.Field(i)

		fieldValue := v.Field(i)

		// tidak di hapus karna ini buat catatan hehehe :D
		// fmt.Printf("Field ke-%d\n", i)
		// fmt.Printf(
		// 	"Field Name: %s\nTag Validator : %s\nField Value: %v\nTipe Data: %s\n\n",
		// 	fieldInfo.Name,
		// 	fieldInfo.Tag.Get("validator"),
		// 	fieldValue.Interface(),
		// 	fieldInfo.Type,
		// )

		validatorSlice[fieldInfo.Name] = validate{
			fieldvalue: fieldValue.Interface(),
			validator:  fieldInfo.Tag.Get("validator"),
		}

	}

	for key, val := range validatorSlice {
		var valueString = fmt.Sprintf("%v", val.fieldvalue)

		if strings.Contains(val.validator, "required") && strings.TrimSpace(valueString) == "" {
			validatorMessage[key] = fmt.Sprintf("Field %s tidak boleh kosong", key)
		}

		if strings.Contains(val.validator, "required") && strings.Contains(val.validator, "number") {
			var valueInt = val.fieldvalue.(int)

			if valueInt == 0 {
				validatorMessage[key] = fmt.Sprintf("Field %s tidak boleh berisi angka 0", key)
			}
		}
	}

	if len(validatorMessage) == 0 {
		return false, validatorMessage
	}

	return true, validatorMessage
}
