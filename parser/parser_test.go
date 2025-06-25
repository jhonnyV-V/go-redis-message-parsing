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

	for _, tc := range testCases {
		actual, err := parser.Deserialise(tc.input)
		if err != nil {
			t.Errorf("deserialise(%q) returned error: %v", tc.input, err)
		} else if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("deserialise(%q) = %v, want %v", tc.input, actual, tc.expected)
		}
	}
}

// func TestSerialise(t *testing.T) {
// }
