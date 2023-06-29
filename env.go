package butler

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	InvalidFormatError = errors.New("please provide a valid time format")

	timeLayouts = []string{
		time.RFC3339,
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,

		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		"2006/01/02T15:04:05",
	}
)

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func LookUpEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// LoadEnvFile loads environment variable key-value pairs from an environment file and sets them.
func LoadEnvFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		pair := strings.SplitN(line, "=", 2)
		if len(pair) != 2 {
			return errors.New("invalid format in the environment file")
		}
		key := strings.ToLower(strings.TrimSpace(pair[0]))
		value := strings.TrimSpace(pair[1])
		if err = SetEnv(key, value); err != nil {
			return err
		}
	}

	return nil
}

// LoadConfig loads configuration data from a file into the provided environment struct.
// The function supports different file types such as .env and .yaml.
func LoadConfig(envStruct interface{}, filePath string) error {
	if envStruct == nil {
		return errors.New("struct cannot be nil")
	}
	if filePath == "" {
		return errors.New("provide file path")
	}

	fileExtension := strings.Split(filePath, ".")
	fileType := fileExtension[len(fileExtension)-1]
	switch fileType {
	case envFileType:
		err := LoadEnvFile(filePath)
		if err != nil {
			return err
		}
		return bind(envStruct)
	case yamlFileType:
		return LoadYAMLFile(envStruct, filePath)
	}

	return nil
}

// bind binds environment variable values to the corresponding fields in the provided environment struct.
func bind(envStruct interface{}) error {
	val := reflect.ValueOf(envStruct)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct || val.IsNil() {
		return errors.New("struct must be a pointer to a struct")
	}
	v := val.Elem()
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(envTag)
		if tag == "" {
			continue
		}

		field := v.Type().Field(i).Tag.Get(envTag)

		envFieldValue, ok := LookUpEnv(field)
		if !ok || envFieldValue == "" {
			envFieldValue = v.Type().Field(i).Tag.Get(defaultTag)
			if envFieldValue == "" {
				continue
			}
		}

		currentFieldValue := reflect.Indirect(v).Field(i).Interface()

		switch currentFieldValue.(type) {
		case string:
			v.Field(i).SetString(envFieldValue)
		case int, int8, int16, int32, int64:
			integer, err := strconv.ParseInt(envFieldValue, 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetInt(integer)
		case uint, uint8, uint16, uint32, uint64:
			integer, err := strconv.ParseUint(envFieldValue, 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetUint(integer)
		case float32, float64:
			float, err := strconv.ParseFloat(envFieldValue, 64)
			if err != nil {
				panic(err)
			}
			v.Field(i).SetFloat(float)
		case bool:
			boolean, err := strconv.ParseBool(envFieldValue)
			if err != nil {
				return err
			}
			v.Field(i).SetBool(boolean)

		case time.Time:
			format := v.Type().Field(i).Tag.Get(formatTag)
			parsedTime, err := parseTimeWithMultipleLayouts(format, envFieldValue)
			if err != nil {
				return err
			}
			// check if it is a pointer
			if _, assertOk := currentFieldValue.(*time.Time); assertOk {
				v.Field(i).Set(reflect.ValueOf(&parsedTime))
			} else {
				v.Field(i).Set(reflect.ValueOf(parsedTime))
			}

		case time.Duration: // 1h0m0s
			d, err := time.ParseDuration(envFieldValue)
			if err != nil {
				return err
			}
			v.Field(i).Set(reflect.ValueOf(d))
		}
	}
	return nil
}

func parseTime(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// parseTimeWithMultipleLayouts parses time with different possible layouts
func parseTimeWithMultipleLayouts(layout, value string) (parsedTime time.Time, err error) {
	if layout != "" {
		parsedTime, err = parseTime(layout, value)
	} else {
		for _, layout = range timeLayouts {
			parsedTime, err = parseTime(layout, value)
			if err == nil {
				break
			}
		}
	}

	return parsedTime, err
}
