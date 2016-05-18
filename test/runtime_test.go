package test

import (
	"testing"
	tools "../"
	"runtime"
	"strings"
	"reflect"
)

func TestFuncLineinfo(t *testing.T) {
	fileName, funcName, lineNum := tools.FuncLineinfo(runtime.Caller(0))

	if !strings.HasSuffix(fileName, "runtime_test.go") {
		t.Error("failed to get the caller file name. The result is",fileName,", but the expect result should contains 'runtimetest.go'")
	}

	if !reflect.DeepEqual(funcName, "TestFuncLineinfo") {
		t.Error("failed to get the caller function name. The result is",funcName,", but the expect result should contains 'TestFuncLineinfo'")
	}

	if lineNum<=0 {
		t.Error("failed to get the caller line number. The result is",lineNum,", but the expect result should great than zero")
	}
}
