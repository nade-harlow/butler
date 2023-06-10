package main

import (
	"bufio"
	"os"
	"strings"
)

type data struct {
	id   string `env:"id"`
	age  int    `env:"age"`
	name string `env:"name"`
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
