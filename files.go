package tools

import (
	"os"
	"path/filepath"
)

func DirWalkUp(path string, walkFun func(string)error) {
	for path != "" && path != "." && path != ".." {
		walkFun(path)

		switch checkPath := filepath.Dir(path); {
		case checkPath == path || path == "." || path == "..":
			return
		case checkPath != "":
			path = checkPath
		}
	}
}

func FilesRemoveAll(path string) {
	DirWalkUp(path, os.RemoveAll)
}
