package populate

import (
	"fmt"
	"github.com/mtabini/go-bowtie"
	"net/http"
	"reflect"
	"strings"
)

func populateArray(name string, sourceType reflect.Type, sourceValue reflect.Value, destType reflect.Type, destValue reflect.Value) error {
	itemCount := sourceValue.Len()

	rowDestType := destType.Elem()

	for index := 0; index < itemCount; index++ {
		rowSourceValue := sourceValue.Index(index)
		rowSourceType := rowSourceValue.Type()

		rowDestValue := reflect.New(rowDestType).Elem()

		if err := populateValue(fmt.Sprintf("%s[%d]", name, index), rowSourceType, rowSourceValue, rowDestType, rowDestValue); err != nil {
			return err
		}

		reflect.Append(destValue, rowDestValue)
	}

	return nil
}

func analyzeField(field reflect.StructField) (name string, ignore bool) {
	if field.Anonymous {
		ignore = true
		return
	}

	name = field.Name

	tag := strings.TrimSpace(field.Tag.Get("json"))

	if tag != "" {
		parts := strings.Split(tag, ",")

		name = strings.TrimSpace(parts[0])

		if name == "-" {
			ignore = true
			return
		}
	}

	tag = strings.TrimSpace(field.Tag.Get("pop"))

	switch tag {
	case "required", "optional":
		break

	default:
		ignore = true
	}

	return
}

func valueForField(sourceType reflect.Type, sourceValue reflect.Value, fieldName string) (value reflect.Value, found bool) {
	switch sourceType.Kind() {
	case reflect.Struct:
		value = sourceValue.FieldByName(fieldName)

	case reflect.Map:
		value = sourceValue.MapIndex(reflect.ValueOf(fieldName))

	default:
		panic("Unexpectedly tried to extract a field from something that is neither a struct nor a map!")
	}

	found = value.IsValid()

	return
}

func copyValue(sourceValue reflect.Value, destValue reflect.Value) error {
	return nil
}

func populateStruct(name string, sourceType reflect.Type, sourceValue reflect.Value, destType reflect.Type, destValue reflect.Value) error {
	fieldCount := destType.NumField()

	for index := 0; index < fieldCount; index++ {
		field := destType.Field(index)
		fieldDestValue := destValue.Field(index)

		if !fieldDestValue.CanSet() {
			continue
		}

		fieldName, ignore := analyzeField(field)

		if ignore {
			continue
		}

		fieldSourceValue, found := valueForField(sourceType, sourceValue, fieldName)

		if found {
			if err := populateValue(name+"."+fieldName, fieldSourceValue.Type(), fieldSourceValue, fieldDestValue.Type(), fieldDestValue); err != nil {
				return err
			}
		}
	}

	return nil
}

func populateValue(name string, sourceType reflect.Type, sourceValue reflect.Value, destType reflect.Type, destValue reflect.Value) error {
	if sourceType.Kind() == reflect.Interface {
		sourceValue = sourceValue.Elem()
		sourceType = sourceValue.Type()
	}

	if sourceType.Kind() == reflect.Ptr {
		if sourceValue.IsNil() {
			return nil
		}

		sourceValue = sourceValue.Elem()
		sourceType = sourceValue.Type()
	}

	if destType.Kind() == reflect.Ptr {
		if destValue.IsNil() {
			destValue.Set(reflect.New(destType.Elem()))
		}

		destValue = destValue.Elem()
		destType = destValue.Type()
	}

	switch destType.Kind() {
	case reflect.Array, reflect.Slice:
		if sourceType.Kind() != reflect.Array {
			return bowtie.NewError(http.StatusBadRequest, "%s.type.invalid", name)
		}

		return populateArray(name, sourceType, sourceValue, destType, destValue)

	case reflect.Struct:
		switch sourceType.Kind() {
		case reflect.Struct, reflect.Map:
			return populateStruct(name, sourceType, sourceValue, destType, destValue)

		default:
			return bowtie.NewError(http.StatusBadRequest, "%s.type.invalid", name)
		}

	default:
		if !sourceValue.Type().ConvertibleTo(destType) {
			return bowtie.NewError(http.StatusBadRequest, "%s.type.invalid", name)
		}

		if destType.Kind() == reflect.String && sourceType.Kind() == reflect.Int {
			return bowtie.NewError(http.StatusBadRequest, "%s.type.invalid", name)
		} else {
			destValue.Set(sourceValue.Convert(destType))
		}

	}

	return nil
}

func Populate(src interface{}, dest interface{}) error {
	sourceValue := reflect.ValueOf(src)
	sourceType := reflect.TypeOf(src)

	destValue := reflect.ValueOf(dest)

	if destValue.Kind() != reflect.Ptr {
		return bowtie.NewError(http.StatusInternalServerError, "`dest` must be a pointer")
	}

	if destValue.IsNil() {
		return bowtie.NewError(http.StatusInternalServerError, "`dest` cannot be nil")
	}

	destValue = destValue.Elem()
	destType := destValue.Type()

	return populateValue("root", sourceType, sourceValue, destType, destValue)
}
