package butler

import (
	"bufio"
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Port struct {
	Number int `env:"number"`
}
type data struct {
	Port    Port   `env:"port"`
	Env     string `env:"env"`
	Verbose bool   `env:"verbose"`
}

const (
	envTag       = "env"
	yamlFileType = "yaml"
	envFileType  = "env"
)

func set(key, value string) error {
	return os.Setenv(key, value)
}

func get(key string) string {
	return os.Getenv(key)
}

func lookUp(key string) (string, bool) {
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
		if err = set(strings.ToLower(pair[0]), pair[1]); err != nil {
			return err
		}
	}

	return nil
}

func SetupConfig(envStruct interface{}, path string) error {
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
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(get(field))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			integer, err := strconv.ParseInt(get(field), 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetInt(integer)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			integer, err := strconv.ParseUint(get(field), 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetUint(integer)
		case reflect.Float32, reflect.Float64:
			float, err := strconv.ParseFloat(get(field), 64)
			if err != nil {
				panic(err)
			}
			v.Field(i).SetFloat(float)
		case reflect.Bool:
			boolean, err := strconv.ParseBool(get(field))
			if err != nil {
				panic(err)
			}
			v.Field(i).SetBool(boolean)
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
