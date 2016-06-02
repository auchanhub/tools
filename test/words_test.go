package test

import (
	"testing"

	tools "./.."
	"strings"
	"reflect"
)

func TestSplitWords(t *testing.T) {
	result := strings.FieldsFunc("hello, word! guys%oppa", tools.SplitWords)
	expected := []string{"hello", "word", "guys", "oppa"}

	if !reflect.DeepEqual(result, expected) {
		t.Error("the result of slip words is", result, ", but the expected result is", expected)
	}
}