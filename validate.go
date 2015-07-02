package popular

import (
	"fmt"
	"github.com/mtabini/go-bowtie"
	"net/http"
	"reflect"
)

func validateStruct(sourceType reflect.Type, sourceValue reflect.Value, name string) error {
	fieldCount := sourceType.NumField()

	for index := 0; index < fieldCount; index++ {
		field := sourceType.Field(index)

		tag := field.Tag.Get("pop")

		if tag == "required" {
			fieldValue := sourceValue.Field(index)
			fieldType := fieldValue.Type()

			wasPtr := false

			if fieldType.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					return bowtie.NewError(http.StatusBadRequest, "%s.%s.required", name, field.Name)
				}

				fieldValue = fieldValue.Elem()
				fieldType = fieldValue.Type()

				wasPtr = true
			}

			switch fieldType.Kind() {
			case reflect.Struct, reflect.Array, reflect.Slice:
				return validate(fieldType, fieldValue, fmt.Sprintf("%s.%s", name, field.Name))

			case reflect.Map:
				if fieldValue.NumField() == 0 {
					return bowtie.NewError(http.StatusBadRequest, "%s.%s.required", name, field.Name)
				}

			default:
				if !wasPtr && fieldType.Comparable() && fieldValue.Interface() == reflect.Zero(fieldType).Interface() {
					return bowtie.NewError(http.StatusBadRequest, "%s.%s.required", name, field.Name)
				}
			}
		}
	}

	return nil
}

func validateArray(sourceType reflect.Type, sourceValue reflect.Value, name string) error {
	itemCount := sourceValue.Len()

	for index := 0; index < itemCount; index++ {
		rowValue := sourceValue.Index(index)
		rowType := sourceValue.Type()

		if err := validate(rowType, rowValue, fmt.Sprintf("%s.%d", name, index)); err != nil {
			return err
		}
	}

	return nil
}

func validate(sourceType reflect.Type, sourceValue reflect.Value, name string) error {
	if sourceType.Kind() == reflect.Ptr {
		if sourceValue.IsNil() {
			return bowtie.NewError(http.StatusInternalServerError, "%s.null", name)
		}

		sourceValue = sourceValue.Elem()
		sourceType = sourceValue.Type()
	}

	switch sourceType.Kind() {
	case reflect.Struct:
		return validateStruct(sourceType, sourceValue, name)

	case reflect.Array, reflect.Slice:
		return validateArray(sourceType, sourceValue, name)

	default:
		return bowtie.NewError(http.StatusInternalServerError, "%s.invalid", name)
	}
}

func Validate(src interface{}, name string) error {
	sourceType := reflect.TypeOf(src)
	sourceValue := reflect.ValueOf(src)

	return validate(sourceType, sourceValue, name)
}
