package butler

import (
	"bufio"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
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
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(GetEnv(field))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			integer, err := strconv.ParseInt(GetEnv(field), 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetInt(integer)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			integer, err := strconv.ParseUint(GetEnv(field), 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetUint(integer)
		case reflect.Float32, reflect.Float64:
			float, err := strconv.ParseFloat(GetEnv(field), 64)
			if err != nil {
				panic(err)
			}
			v.Field(i).SetFloat(float)
		case reflect.Bool:
			boolean, err := strconv.ParseBool(GetEnv(field))
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

//get fist key
//get value
//get first key and nested value

func yamlReader(envStruct interface{}, path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	s := bufio.NewScanner(f)
	m := make(map[string]interface{})
	key := ""
	parentKey := ""
	value := ""

	for s.Scan() {
		line := s.Text()
		keyValue := strings.Split(line, ":")
		var tempSlice []string
		if keyValue[1] == "" {
			tempSlice = append(tempSlice, keyValue[0])
			keyValue = tempSlice
		}
		key = keyValue[0]
		if len(keyValue) == 1 {
			parentKey = keyValue[0]
			m[keyValue[0]] = strings.TrimSpace(value)
			continue
		}
		if len(keyValue) > 1 {
			value = keyValue[1]
			if strings.HasPrefix(key, " ") {
				appender(parentKey, keyValue, m)
			} else {
				m[key] = strings.TrimSpace(value)
			}
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &envStruct)
	if err != nil {
		log.Fatal(err)
	}

}

var subMap = make(map[string]interface{})

func appender(k string, line []string, v map[string]interface{}) map[string]interface{} {
	subMap[strings.TrimSpace(line[0])] = strings.TrimSpace(line[1])
	v[k] = subMap
	return v
}
