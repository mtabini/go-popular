package popular

import (
	"testing"
)

// type testStruct struct {
// 	Field1 *string `json:"f1" pop:"required"`
// 	Field2 *int    `pop:"required"`
// 	Field3 string  `pop:"optional"`
// 	Field4 int
// 	field5 int
// }

// type testStruct2 struct {
// 	NewField1 string     `pop:"required"`
// 	NewField2 testStruct `pop:"required"`
// }

func TestValidate(t *testing.T) {
	s := &testStruct{}

	err := Validate(s, "test")

	if err == nil || err.Error() != "test.Field1.required" {
		t.Errorf("Unexpected error: %s", err)
	}

	ss := testStruct{}

	err = Validate(ss, "test")

	if err == nil || err.Error() != "test.Field1.required" {
		t.Errorf("Unexpected error: %s", err)
	}

	ss = testStruct{}

	str := "test"
	ss.Field1 = &str

	err = Validate(ss, "test")

	if err == nil || err.Error() != "test.Field2.required" {
		t.Errorf("Unexpected error: %s", err)
	}

	i := 0
	ss.Field2 = &i

	err = Validate(ss, "test")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestValidateNested(t *testing.T) {
	s := testStruct2{}

	err := Validate(s, "test")

	if err == nil || err.Error() != "test.NewField1.required" {
		t.Errorf("Unexpected error: %s", err)
	}

	s.NewField1 = "test"

	err = Validate(s, "test")

	if err == nil || err.Error() != "test.NewField2.Field1.required" {
		t.Errorf("Unexpected error: %s", err)
	}

	str := "test"
	s.NewField2.Field1 = &str

	err = Validate(s, "test")

	if err == nil || err.Error() != "test.NewField2.Field2.required" {
		t.Errorf("Unexpected error: %s", err)
	}

	i := 0
	s.NewField2.Field2 = &i

	err = Validate(s, "test")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
