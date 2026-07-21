package formutil

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type validate struct {
	fieldvalue interface{}
	validator  string
	inputName  string
}

func Validate(reqs any, r *http.Request) (bool, map[string]string) {
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

		validatorSlice[fieldInfo.Name] = validate{
			fieldvalue: fieldValue.Interface(),
			validator:  fieldInfo.Tag.Get("validator"),
			inputName:  fieldInfo.Tag.Get("input_name"),
		}

	}

	for _, val := range validatorSlice {
		var valueString = fmt.Sprintf("%v", val.fieldvalue)

		if strings.Contains(val.validator, "required") {

			if strings.TrimSpace(valueString) == "" {
				validatorMessage[val.inputName] = fmt.Sprintf("Field %s tidak boleh kosong", val.inputName)
			} else if strings.Contains(val.validator, "number") {
				var valueInt = val.fieldvalue.(int)

				if valueInt == 0 {
					validatorMessage[val.inputName] = fmt.Sprintf("Field %s tidak boleh berisi angka 0", val.inputName)
				}

			} else if strings.Contains(val.validator, "match") {
				re := regexp.MustCompile("(?i)match\\[[^\\]]*\\]")
				var indexStart = re.FindStringIndex(val.validator)[0]
				var matchValidator = val.validator[indexStart:]

				replacer := strings.NewReplacer("match", "", "[", "", "]", "")
				sliceInput := strings.Split(replacerMatch(replacer, matchValidator), ".")

				if r.FormValue(sliceInput[0]) != r.FormValue(sliceInput[1]) {
					validatorMessage[val.inputName] = fmt.Sprintf("Field %s harus sama dengan kolom konfirmasi password", val.inputName)
				}

			}

		}

	}

	if len(validatorMessage) == 0 {
		return false, validatorMessage
	}

	return true, validatorMessage
}

func replacerMatch(s *strings.Replacer, str string) string {
	return s.Replace(str)
}
