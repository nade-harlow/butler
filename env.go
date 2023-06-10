package main

import (
	"bufio"
	"errors"
	"fmt"
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

//read from env
//set env
//update env
//read env and bind to struct
//read from file and bind to struct (.yaml and .env)

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
		if strings.Contains(line, "#") {
			if strings.HasPrefix(line, "#") {
				continue
			}
		}
		pair := strings.Split(line, "=")
		if err = set(strings.ToLower(pair[0]), pair[1]); err != nil {
			return err
		}
	}

	return nil
}

func env(envStruct interface{}, path string) error {
	if envStruct == nil {
		return errors.New("struct cannot be nil")
	}
	if path == "" {
		return errors.New("provide file path")
	}
	fileExtension := strings.Split(path, ".")
	fileExt := fileExtension[len(fileExtension)-1]
	switch fileExt {
	case "env":
		err := loadENVFile(path)
		if err != nil {
			return err
		}
		return bind(envStruct)
	case "yaml":
		return loadYAMLFile(envStruct, path)

	}

	return nil
}

func bind(envStruct interface{}) error {
	val := reflect.ValueOf(envStruct)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct || val.IsNil() {
		return fmt.Errorf("struct must be a pointer to a struct")
	}
	v := val.Elem()
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get("env")
		if tag == "" {
			continue
		}

		field := v.Type().Field(i).Tag.Get("env")
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

	err = yaml.Unmarshal(f, envStruct)
	if err != nil {
		return err
	}

	return nil
}
