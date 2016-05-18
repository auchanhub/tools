package tools

import (
	"runtime"
	"strings"
	"path/filepath"
)

func FuncLineinfo(pc uintptr, file string, line int, _ bool) (fileName string, funcName string, lineNum int) {
	fileName = file
	funcName = "unknown"
	lineNum = line

	pc = uintptr(pc) - 1
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return
	}

	fileName, lineNum = fn.FileLine(pc)

	funcName = fn.Name()
	packageName := ""
	if index := strings.LastIndex(funcName, "."); index>=0 {
		packageName = funcName[:index]
		funcName = funcName[index+1:]
	}

	if index := strings.LastIndex(fileName, string(filepath.Separator)); index>=0 {
		fileName = filepath.Join(packageName, fileName[index+1:])
	}

	return
}
