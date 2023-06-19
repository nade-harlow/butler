package butler

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	envTag       = "env"
	yamlFileType = "yaml"
	envFileType  = "env"
)

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func lookUpEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func loadENVFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		if line == "" {
			continue
		}
		if strings.Contains(line, "#") && strings.HasPrefix(line, "#") {
			continue
		}
		pair := strings.Split(line, "=")
		key := strings.ToLower(strings.TrimSpace(pair[0]))
		value := strings.TrimSpace(pair[1])
		if err = SetEnv(key, value); err != nil {
			return err
		}
	}

	return nil
}

func SetupEvn(envStruct interface{}, path string) error {
	if envStruct == nil {
		return errors.New("struct cannot be nil")
	}
	if path == "" {
		return errors.New("provide file path")
	}
	fileExtension := strings.Split(path, ".")
	fileType := fileExtension[len(fileExtension)-1]
	switch fileType {
	case envFileType:
		err := loadENVFile(path)
		if err != nil {
			return err
		}
		return bind(envStruct)
	case yamlFileType:
		return loadYAMLFile(envStruct, path)
	}

	return nil
}

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

		envFieldValue := GetEnv(field)
		if envFieldValue == "" {
			continue
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
			if envValue, ok := lookUpEnv(envTag); ok && envValue != "" {
				format := v.Type().Field(i).Tag.Get("format")
				if format == "" {
					format = "2006-01-02T15:04:05"
				}
				t, err := time.Parse(format, envFieldValue)
				if err != nil {
					return err
				}
				// check if it is a pointer
				if _, ok := currentFieldValue.(*time.Time); ok {
					v.Field(i).Set(reflect.ValueOf(&t))
				} else {
					v.Field(i).Set(reflect.ValueOf(t))
				}
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

func loadYAMLFile(envStruct interface{}, filepath string) error {
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(f, envStruct)
}
