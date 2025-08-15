package parser_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jhonnyV-V/redis-clone/parser"
)

func TestDeserialise(t *testing.T) {
	//cases:
	// "$-1\r\n" nil
	// "*1\r\n$4\r\nping\r\n” []string{"ping"}
	// "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n” []string{"echo", "hello world"}
	// "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n” []string{"get", "key"}
	// "+OK\r\n" "OK"
	// "-Error message\r\n" fmt.Errorf("Error message")
	// "$0\r\n\r\n" []string{}
	// "+hello world\r\n" "hello world"

	//TODO: after migrating to generics, use arrays of multiple types
	testCases := []struct {
		input    string
		expected interface{}
	}{
		{"$-1\r\n", nil},
		{"*1\r\n$4\r\nping\r\n", []string{"ping"}},
		{"*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n", []string{"echo", "hello world"}},
		{"*2\r\n$3\r\nget\r\n$3\r\nkey\r\n", []string{"get", "key"}},
		{"+OK\r\n", "OK"},
		{"-Error message\r\n", fmt.Errorf("Error message")},
		{"$0\r\n\r\n", ""},
		{"$1\r\nh\r\n", "h"},
		{"$2\r\nhi\r\n", "hi"},
		{"+hello world\r\n", "hello world"},
		{":-10\r\n", -10},
		{":+10\r\n", 10},
		{":20\r\n", 20},
		{"#t\r\n", true},
		{"#f\r\n", false},
		{"_\r\n", nil},
	}

	//TODO: Use generics instead of interface{}
	for _, tc := range testCases {
		var i int = 0
		actual, err := parser.Deserialise(tc.input, &i)
		if err != nil {
			t.Errorf("deserialise(%q) returned error: %v", tc.input, err)
		} else if !reflect.DeepEqual(actual, tc.expected) {
			expectedStringSlice, ok := tc.expected.([]string)
			if ok {
				actualStringSlice, err := convertInterfaceSliceToStringSlice(actual.([]interface{}))
				if err == nil {
					if !reflect.DeepEqual(actualStringSlice, expectedStringSlice) {
						t.Errorf("deserialise(%q) = %v, want %v", tc.input, actualStringSlice, expectedStringSlice)
					}
				}
			} else {
				t.Errorf("deserialise(%q) = %v, want %v", tc.input, actual, tc.expected)
			}
		}
	}
}

// TODO: maybe make this function generic to test each individual item
func convertInterfaceSliceToStringSlice(input []interface{}) ([]string, error) {
	result := make([]string, len(input))
	for i, v := range input {
		if str, ok := v.(string); ok {
			result[i] = str
		} else {
			return nil, fmt.Errorf("element at index %d is not a string: %T (value: %v)", i, v, v)
		}
	}
	return result, nil
}

// func TestSerialise(t *testing.T) {
// }
