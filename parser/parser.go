package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func isCRLF(raw string, i int) bool {
	if raw[i] == '\r' && raw[i+1] == '\n' {
		return true
	}
	return false
}

func deserialiseSimpleString(raw string) (string, error) {
	value := strings.Builder{} 

	for i := 1; i < len(raw); i++ {
		if isCRLF(raw, i) {
			return value.String(), nil
		}
		value.WriteByte(raw[i])
	}

	return "", fmt.Errorf("simple string not terminated")
}

func deserialiseErrorString(raw string) (error, error) {
	value := strings.Builder{} 

	for i := 1; i < len(raw); i++ {
		if isCRLF(raw, i) {
			return errors.New(value.String()), nil
		}
		value.WriteByte(raw[i])
	}

	return nil, fmt.Errorf("error not terminated")
}

func deserialiseBulkString(raw string) (interface{}, error)  {
	toParse := string(raw[1:])
	valuesToParse := 0
	rawValues := ""

	if (toParse[0] == '-') {
		return nil, nil
	}

	for i := 0; i < len(toParse); i++ {
		if isCRLF(toParse, i) {
			if rawValues == "" || rawValues == "0" {
				return "", nil
			}
			toParse = toParse[i+2:]
			break
		}
		rawValues += string(toParse[i])
	}

	valuesToParse, err := strconv.Atoi(rawValues)
	if err != nil {
		return "", fmt.Errorf("Failed to parse values (%q) err: %w", rawValues, err)
	}
	
	return toParse[:valuesToParse], nil
}

func Deserialise(raw string) (interface{}, error)  {
	switch raw[0] {
	case '+':
		return deserialiseSimpleString(raw)
	case '-':
		return deserialiseErrorString(raw)
	case '$':
		return deserialiseBulkString(raw)
	case '*':
		return nil, nil
	}

	return "", fmt.Errorf("unknown type: %#v", raw[0])
}

func Serialise(value string) string {
	return ""
}
