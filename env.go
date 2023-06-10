package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type data struct {
	Port float64 `env:"PORT"`
	Env  string  `env:"ENV"`
	//Name string `env:"name"`
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

func load(path string) error {
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
		if err = set(pair[0], pair[1]); err != nil {
			return err
		}
	}

	return nil
}

func env(envStruct interface{}, path string) error {
	if path == "" {
		return errors.New("provide file path")
	}
	err := load(path)
	if err != nil {
		return err
	}
	if envStruct == nil {
		return errors.New("struct cannot be nil")
	}

	return bind(envStruct)
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
		}
	}
	return nil
}
