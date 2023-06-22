package butler

import (
	"reflect"
	"testing"
)

func TestLoadYAMLFile(t *testing.T) {
	yamlContent := yamlFileContent()

	file, err := createTempFile(yamlContent, "config.*.yaml")
	if err != nil {
		t.Errorf("createTempFile returned an error: %s", err.Error())
	}

	envStruct := &TestStruct{}
	err = LoadYAMLFile(&envStruct, file.Name())
	if err != nil {
		t.Errorf("LoadYAMLFile returned an error: %s", err.Error())
	}

	// Verify the loaded values
	expectedStruct := &TestStruct{
		StringField: "test value",
		IntField:    123,
		UintField:   456,
		FloatField:  3.14,
		BoolField:   true,
	}

	if !reflect.DeepEqual(envStruct, expectedStruct) {
		t.Errorf("Loaded struct does not match the expected values:\nExpected: %+v\nGot: %+v", expectedStruct, envStruct)
	}
}

func TestRemoveHyphen(t *testing.T) {
	key := "test_key"
	expected := "testkey"

	result := removeHyphen(key)
	if result != expected {
		t.Errorf("removeHyphen returned an unexpected result: expected '%s', got '%s'", expected, result)
	}
}

func TestGetValueWithType(t *testing.T) {
	t.Run("Test boolean", func(t *testing.T) {
		boolInput := "true"
		boolExpected := true
		boolResult := getValueWithType(boolInput)
		if boolResult != boolExpected {
			t.Errorf("getValueWithType returned an unexpected result for boolean input: expected '%v', got '%v'", boolExpected, boolResult)
		}
	})

	t.Run("Test integer", func(t *testing.T) {
		intInput := "-123"
		intExpected := int64(-123)
		intResult := getValueWithType(intInput)
		if intResult != intExpected {
			t.Errorf("getValueWithType returned an unexpected result for integer input: expected '%v', got '%v'", intExpected, intResult)
		}
	})

	t.Run("Test float", func(t *testing.T) {
		floatInput := "3.14"
		floatExpected := 3.14
		floatResult := getValueWithType(floatInput)
		if floatResult != floatExpected {
			t.Errorf("getValueWithType returned an unexpected result for float input: expected '%v', got '%v'", floatExpected, floatResult)
		}
	})

	t.Run("Test uint", func(t *testing.T) {
		unitInput := "456"
		unitExpected := uint64(456)
		unitResult := getValueWithType(unitInput)
		if unitResult != unitExpected {
			t.Errorf("getValueWithType returned an unexpected result for unit input: expected '%v', got '%v'", unitExpected, unitResult)
		}
	})

	t.Run("Test string", func(t *testing.T) {
		stringInput := "test"
		stringExpected := "test"
		stringResult := getValueWithType(stringInput)
		if stringResult != stringExpected {
			t.Errorf("getValueWithType returned an unexpected result for string input: expected '%v', got '%v'", stringExpected, stringResult)
		}
	})

}

func TestAppender(t *testing.T) {
	parentKey := "parent"
	line := []string{"child_key", "child_value"}

	v := make(map[string]interface{})
	expected := map[string]interface{}{
		parentKey: map[string]interface{}{
			"childkey": "child_value",
		},
	}

	result := appender(parentKey, line, v)

	// Check if the modified map matches the expected result
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("appender returned an unexpected result:\nExpected: %+v\nActual: %+v", expected, result)
	}

	// Check if the subMap is modified correctly
	if subMap["childkey"] != "child_value" {
		t.Errorf("appender did not set the child value correctly")
	}
}
