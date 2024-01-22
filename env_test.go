package butler

import (
	"bufio"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSetEnv(t *testing.T) {
	key := "TEST_KEY"
	value := "test value"

	err := SetEnv(key, value)
	if err != nil {
		t.Errorf("SetEnv returned an error: %s", err.Error())
	}

	envValue := os.Getenv(key)
	if envValue != value {
		t.Errorf("Expected environment variable value: %s, got: %s", value, envValue)
	}
}

func TestGetEnv(t *testing.T) {
	key := "TEST_KEY"
	value := "test value"
	os.Setenv(key, value)

	envValue := GetEnv(key)
	if envValue != value {
		t.Errorf("Expected environment variable value: %s, got: %s", value, envValue)
	}
}

func TestLookUpEnv(t *testing.T) {
	key := "TEST_KEY"
	value := "test value"
	os.Setenv(key, value)

	envValue, ok := LookUpEnv(key)
	if !ok {
		t.Errorf("Expected to find environment variable with key: %s", key)
	}
	if envValue != value {
		t.Errorf("Expected environment variable value: %s, got: %s", value, envValue)
	}
}

func TestLoadEnvFile(t *testing.T) {
	// Create a temporary environment file for testing
	fileContent := `# Comment line
	ENV_VAR1=Value1
	ENV_VAR2=Value2`
	tempFile, err := createTempFile(fileContent, "envfile.*.env")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	err = LoadEnvFile(tempFile.Name())
	if err != nil {
		t.Errorf("LoadEnvFile returned an error: %v", err)
	}

	expectedCalls := [][]string{
		{"ENV_VAR1", "Value1"},
		{"ENV_VAR2", "Value2"},
	}
	for i, call := range expectedCalls {
		val := GetEnv(call[0])
		if val != call[1] {
			t.Errorf("SetEnv call %d does not match the expected value. Got: %v, Expected: %v", i+1, call, expectedCalls[i])
		}
	}
}

type TestStruct struct {
	StringField   string        `env:"STRING_FIELD"`
	IntField      int           `env:"INT_FIELD"`
	UintField     uint          `env:"UINT_FIELD"`
	FloatField    float64       `env:"FLOAT_FIELD"`
	BoolField     bool          `env:"BOOL_FIELD"`
	TimeField     time.Time     `env:"TIME_FIELD" format:"2006-01-02"`
	DurationField time.Duration `env:"DURATION_FIELD"`
}

func TestBind(t *testing.T) {
	setEnv(t)
	testStruct := &TestStruct{}
	err := bind(testStruct)
	if err != nil {
		t.Errorf("bind returned an error: %v", err)
	}

	// Verify the field values
	expectedStruct := &TestStruct{
		StringField:   "test value",
		IntField:      123,
		UintField:     456,
		FloatField:    3.14,
		BoolField:     true,
		TimeField:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		DurationField: time.Hour + 30*time.Minute,
	}

	if !reflect.DeepEqual(testStruct, expectedStruct) {
		t.Errorf("Binded struct does not match the expected value.\nGot: %#v\nExpected: %#v", testStruct, expectedStruct)
	}
}

func TestLoadConfig(t *testing.T) {
	setEnv(t)

	// Verify the field values
	expectedStruct := &TestStruct{
		StringField:   "test value",
		IntField:      123,
		UintField:     456,
		FloatField:    3.14,
		BoolField:     true,
		TimeField:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		DurationField: time.Hour + 30*time.Minute,
	}
	t.Run("Test loading .env file", func(t *testing.T) {
		file, err := createTempFile(envFileContent(), "envfile.*.env")
		if err != nil {
			return
		}
		testStruct := &TestStruct{}
		err = LoadConfig(testStruct, file.Name())
		if err != nil {
			t.Errorf("LoadConfig returned an error: %v", err)
		}

		if !reflect.DeepEqual(testStruct, expectedStruct) {
			t.Errorf("Loaded struct does not match the expected value.\nGot: %#v\nExpected: %#v", testStruct, expectedStruct)
		}
	})
	t.Run("Test Loading .yaml file", func(t *testing.T) {
		file, err := createTempFile(yamlFileContent(), "config.*.yaml")
		if err != nil {
			return
		}
		testStruct := &TestStruct{}
		err = LoadConfig(testStruct, file.Name())
		if err != nil {
			t.Errorf("LoadConfig returned an error: %v", err)
		}
		expectedStruct := &TestStruct{
			StringField: "test value",
			IntField:    123,
			UintField:   456,
			FloatField:  3.14,
			BoolField:   true,
		}

		if !reflect.DeepEqual(testStruct, expectedStruct) {
			t.Errorf("Loaded struct does not match the expected value.\nGot: %#v\nExpected: %#v", testStruct, expectedStruct)
		}
	})

}

func yamlFileContent() string {
	return `STRING_FIELD: test value
INT_FIELD: 123
UINT_FIELD: 456
FLOAT_FIELD: 3.14
BOOL_FIELD: true`
}

func envFileContent() string {
	return `STRING_FIELD=test value
			INT_FIELD=123
			UINT_FIELD=456
			FLOAT_FIELD=3.14
			BOOL_FIELD=true
			TIME_FIELD=2022-01-01
			DURATION_FIELD=1h30m`
}

func setEnv(t *testing.T) {
	// Set up test environment variables
	err := SetEnv("STRING_FIELD", "test value")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	err = SetEnv("INT_FIELD", "123")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	err = SetEnv("UINT_FIELD", "456")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	err = SetEnv("FLOAT_FIELD", "3.14")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	err = SetEnv("BOOL_FIELD", "true")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	err = SetEnv("TIME_FIELD", "2022-01-01")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	err = SetEnv("DURATION_FIELD", "1h30m")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
}

func createTempFile(content, pattern string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(tempFile)
	_, err = writer.WriteString(content)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}
