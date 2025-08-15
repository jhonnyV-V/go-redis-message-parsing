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

func deserialiseSimpleString(raw string, pointer *int) (string, error) {
	value := strings.Builder{}

	*pointer += 1
	for i := *pointer; i < len(raw); i++ {
		*pointer += 1
		if isCRLF(raw, i) {
			return value.String(), nil
		}
		value.WriteByte(raw[i])
	}

	return "", fmt.Errorf("simple string not terminated")
}

func deserialiseErrorString(raw string, pointer *int) (error, error) {
	value := strings.Builder{}

	*pointer += 1
	for i := 1; i < len(raw); i++ {
		*pointer += 1
		if isCRLF(raw, i) {
			*pointer += 1
			return errors.New(value.String()), nil
		}
		value.WriteByte(raw[i])
	}

	return nil, fmt.Errorf("error not terminated")
}

func deserialiseBulkString(raw string, pointer *int) (interface{}, error) {
	toParse := string(raw[1:])
	valuesToParse := 0
	rawValues := ""

	if toParse[0] == '-' {
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

	*pointer += len(rawValues) + 1 + valuesToParse + 2 + 2

	return toParse[:valuesToParse], nil
}

func deserialiseInteger(raw string, pointer *int) (int, error) {
	value := 0
	rawValues := ""

	*pointer += 1
	for i := 1; i < len(raw); i++ {
		*pointer += 1
		if isCRLF(raw, i) {
			*pointer += 1
			break
		}
		rawValues += string(raw[i])
	}

	value, err := strconv.Atoi(rawValues)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse int (%q) err: %w", rawValues, err)
	}

	return value, nil
}

func deserialiseBool(raw string, pointer *int) (bool, error) {
	*pointer += 4
	if raw[1] == 't' {
		return true, nil
	}

	if raw[1] == 'f' {
		return false, nil
	}

	return false, fmt.Errorf("Failed to parse bool (%q)", raw)
}

func deserialiseArray(raw string, pointer *int) ([]interface{}, error) {
	toParse := string(raw[1:])
	valuesToParse := 0
	rawValues := ""
	values := []interface{}{}

	for i := 0; i < len(toParse); i++ {
		if isCRLF(toParse, i) {
			if rawValues == "" || rawValues == "0" {
				return values, nil
			}
			toParse = toParse[i+2:]
			break
		}
		rawValues += string(toParse[i])
	}

	valuesToParse, err := strconv.Atoi(rawValues)
	if err != nil {
		return values, fmt.Errorf("Failed to parse values (%q) err: %w", rawValues, err)
	}

	*pointer += len(rawValues) + 3

	for i := 0; i < valuesToParse; i++ {
		value, err := Deserialise(raw[*pointer:], pointer)
		if err != nil {
			return values, fmt.Errorf("deserialise(%q) returned error: %v", raw, err)
		}
		values = append(values, value)
	}

	return values, nil
}

func Deserialise(raw string, pointer *int) (interface{}, error) {
	switch raw[0] {
	case '+':
		return deserialiseSimpleString(raw, pointer)
	case '-':
		return deserialiseErrorString(raw, pointer)
	case '$':
		return deserialiseBulkString(raw, pointer)
	case ':':
		return deserialiseInteger(raw, pointer)
	case '#':
		return deserialiseBool(raw, pointer)
	case '_': // null
		*pointer += 3
		return nil, nil
	case '*':
		return deserialiseArray(raw, pointer)
	}

	return "", fmt.Errorf("unknown type: %#v", raw[0])
}

func Serialise(value string) string {
	return ""
}
