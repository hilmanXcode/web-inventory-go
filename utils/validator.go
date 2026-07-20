package utils

import (
	"fmt"
	"reflect"
	"strings"
)

type validate struct {
	fieldvalue interface{}
	validator  string
	inputName  string
}

func Validate(reqs any) (bool, map[string]string) {
	var validatorMessage = map[string]string{}

	t := reflect.TypeOf(reqs)
	v := reflect.ValueOf(reqs)

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
			inputName:  fieldInfo.Tag.Get("input_name"),
		}

	}

	for _, val := range validatorSlice {
		var valueString = fmt.Sprintf("%v", val.fieldvalue)

		if strings.Contains(val.validator, "required") && strings.TrimSpace(valueString) == "" {
			validatorMessage[val.inputName] = fmt.Sprintf("Field %s tidak boleh kosong", val.inputName)
		}

		if strings.Contains(val.validator, "required") && strings.Contains(val.validator, "number") {
			var valueInt = val.fieldvalue.(int)

			if valueInt == 0 {
				validatorMessage[val.inputName] = fmt.Sprintf("Field %s tidak boleh berisi angka 0", val.inputName)
			}
		}
	}

	if len(validatorMessage) == 0 {
		return false, validatorMessage
	}

	return true, validatorMessage
}
