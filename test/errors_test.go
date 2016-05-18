package test

import (
	"testing"
	"github.com/pkg/errors"
	"os"
	"github.com/7phones/tools"
	"runtime"
	"strings"
	"reflect"
	"fmt"
)

func TestErrorsDump(t *testing.T) {
	testData := []string{
		"test invalid args",
		"test level 1",
		"test level 2",
	}
	initErr := os.ErrInvalid
	err := initErr

	// generate the test set
	for _, message := range testData {
		err = errors.Wrap(err, message)
	}

	fileName, _, lineNum := tools.FuncLineinfo(runtime.Caller(0))

	// build the check set
	revert := func(data []string) (result []string) {
		for _, s := range data {
			result = append([]string{s}, result...)
		}

		return
	}
	prefix := fmt.Sprintf("%s:%d: ", fileName, lineNum - 3)
	expect := prefix + strings.Join(revert(testData), "\n" + prefix) + "\n" + initErr.Error() + "\n"

	// execute the test function
	exist := tools.ErrorsDump(err)

	if !reflect.DeepEqual(expect, exist) {
		t.Error("failed to dump error. The result is\r\n", exist, "\r\nbut expect result is\r\n", expect)
	}
}