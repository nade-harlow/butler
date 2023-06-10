package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type data struct {
	Port string `env:"PORT"`
	Env  string `env:"ENV"`
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

		switch v.Field(i).Kind() {
		case reflect.String:
			field := v.Type().Field(i).Tag.Get("env")
			v.Field(i).SetString(get(field))
		}
	}
	return nil
}
