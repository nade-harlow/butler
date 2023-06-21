package butler

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
)

// LoadYAMLFile loads environment variable key-value pairs from a YAML file and populates the provided environment struct.
// The function calls the yamlReader function internally to read and process the YAML file.
func LoadYAMLFile(envStruct interface{}, filepath string) error {
	return yamlReader(envStruct, filepath)
}

func yamlReader(envStruct interface{}, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.New("error opening .yaml file: " + err.Error())
	}
	s := bufio.NewScanner(f)
	m := make(map[string]interface{})
	key := ""
	parentKey := ""
	value := ""

	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		keyValue := strings.Split(line, ":")
		var tempSlice []string
		if keyValue[1] == "" {
			tempSlice = append(tempSlice, keyValue[0])
			keyValue = tempSlice
		}
		key = removeHyphen(keyValue[0])
		if len(keyValue) == 1 {
			parentKey = removeHyphen(keyValue[0])
			m[keyValue[0]] = strings.TrimSpace(value)
			continue
		}
		if len(keyValue) > 1 {
			value = keyValue[1]
			if strings.HasPrefix(key, " ") {
				appender(parentKey, keyValue, m)
			} else {
				subMap = make(map[string]interface{}) //reset map
				val := getValueWithType(strings.TrimSpace(value))
				m[key] = val
			}
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return errors.New("error marshaling .yaml: " + err.Error())
	}
	err = json.Unmarshal(b, &envStruct)
	if err != nil {
		return errors.New("error unmarshalling .yaml: " + err.Error())
	}
	return nil
}

var subMap = make(map[string]interface{})

// appender appends a nested key-value pair to the parentKey in the provided map.
func appender(parentKey string, line []string, v map[string]interface{}) map[string]interface{} {
	childValue := getValueWithType(strings.TrimSpace(line[1]))
	childKey := removeHyphen(strings.TrimSpace(line[0]))
	subMap[childKey] = childValue
	v[parentKey] = subMap
	return v
}

// getValueWithType tries to parse the input string into different data types and returns the parsed value.
func getValueWithType(input string) interface{} {

	// Try parsing as unit
	if val, err := strconv.ParseUint(input, 10, 64); err == nil {
		return val
	}
	// Try parsing as boolean
	if val, err := strconv.ParseBool(input); err == nil {
		return val
	}
	// Try parsing as integer
	if val, err := strconv.ParseInt(input, 10, 64); err == nil {
		return val
	}
	// Try parsing as float
	if val, err := strconv.ParseFloat(input, 64); err == nil {
		return val
	}
	// Return as string if no other type matched
	return input
}

func removeHyphen(key string) string {
	return strings.Replace(key, "_", "", -1)
}
