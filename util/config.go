package util

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	commentPrefix   = "-- "
	lineSplit       = "="
	stringSeparator = "\""
	configLUAFile   = "/config.lua"
)

// ConfigLUA holds the parsed config lua file
type ConfigLUA struct {
	v map[string]interface{}
}

// ParseConfig parses the config lua file
func ParseConfig(path string) {
	Config = &ConfigLUA{
		make(map[string]interface{}),
	}
	file, err := os.Open(path + configLUAFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		if strings.HasPrefix(currentLine, commentPrefix) || currentLine == "" {
			continue
		}
		args := strings.Split(currentLine, lineSplit)
		if len(args) != 2 {
			continue
		}
		paramName := strings.TrimSpace(args[0])
		paramValue := strings.TrimSpace(args[1])
		if strings.HasPrefix(paramValue, stringSeparator) {
			paramValue = strings.TrimPrefix(paramValue, stringSeparator)
			paramValue = strings.TrimSuffix(paramValue, stringSeparator)
		}
		Config.v[paramName] = paramValue
	}
}

// String returns a config lua string value
func (c *ConfigLUA) String(key string) string {
	if v, ok := c.v[key].(string); ok {
		return v
	}
	return ""
}

// Int returns a config lua int value
func (c *ConfigLUA) Int(key string) int {
	k, err := strconv.Atoi(c.v[key].(string))
	if err != nil {
		return 0
	}
	return k
}

// Bool returns a config lua bool value
func (c *ConfigLUA) Bool(key string) bool {
	k, err := strconv.ParseBool(c.v[key].(string))
	if err != nil {
		return false
	}
	return k
}
