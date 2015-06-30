package populate

import (
	"testing"
)

type testStruct struct {
	Field1 *string `json:"f1" pop:"required"`
	Field2 *int    `pop:"required"`
	Field3 string  `pop:"optional"`
	Field4 int
	field5 int
}

type testStruct2 struct {
	NewField1 string     `pop:"required"`
	NewField2 testStruct `pop:"required"`
}

func TestPopulateStruct(t *testing.T) {
	s := map[string]interface{}{
		"f1":     "test",
		"Field2": 123,
	}

	d := testStruct{}

	if err := Populate(s, &d); err != nil {
		t.Errorf("%s", err)
	}

	if *d.Field1 != "test" {
		t.Errorf("Field1: Expected `test`, got %#v instead", *d.Field1)
	}

	if *d.Field2 != 123 {
		t.Errorf("Field2: Expected `test`, got %#v instead", *d.Field2)
	}

	s = map[string]interface{}{
		"NewField1": "testing",
		"NewField2": map[string]interface{}{
			"f1":     "test",
			"Field2": 123,
		},
	}

	dd := testStruct2{}

	if err := Populate(s, &dd); err != nil {
		t.Errorf("%s", err)
	}

	if dd.NewField1 != "testing" {
		t.Errorf("NewField1: Expected `testing`, got %#v instead", dd.NewField1)
	}

	if *dd.NewField2.Field1 != "test" {
		t.Errorf("Field1: Expected `test`, got %#v instead", *dd.NewField2.Field1)
	}

	if *dd.NewField2.Field2 != 123 {
		t.Errorf("Field2: Expected `test`, got %#v instead", *dd.NewField2.Field2)
	}
}

func TestPopulateValue(t *testing.T) {
	x := ""

	if err := Populate("123", &x); err != nil {
		t.Errorf("Value: %s", err)
	}

	if x != "123" {
		t.Errorf("Expected `123`, got %s instead", x)
	}
}

func TestPopulateNumericString(t *testing.T) {
	x := 0

	if err := Populate("123", &x); err == nil {
		t.Error("Able to decode a numeric string into a numeric field.")
	}
}

func TestPopulateError(t *testing.T) {
	x := ""

	if err := Populate(123, &x); err == nil {
		t.Errorf("Able to populate string with number: %s", x)
	}
}

func TestConversion(t *testing.T) {
	x := 0.0

	if err := Populate(123, &x); err != nil {
		t.Errorf("Conversion: %s", err)
	}

	if x != 123.0 {
		t.Errorf("Conversion failed. Expected `123.0`, got %f instead", x)
	}
}
