package test

import (
	tools "../"
	"testing"
	"os"
	"reflect"
	"strings"
	"path/filepath"
)

func TestDirWalkUp(t *testing.T) {
	testRootDir := "data"
	testData := []struct {
		dir      string
		expected []string
	}{
		{
			dir: "data/test/path",
			expected: []string{
				"data/test/path",
				"data/test",
				"data",
			},
		},
		{
			dir: "data/test/other.path/пример",
			expected: []string{
				"data/test/other.path/пример",
				"data/test/other.path",
				"data/test",
				"data",
			},
		},
	}

	defer os.RemoveAll(testRootDir)

	for _, test := range testData {
		if err := os.MkdirAll(test.dir, os.ModePerm); err != nil {
			t.Error("filed to create the test directory", test.dir, ":", err)

			continue
		}

		exist := []string{}

		tools.DirWalkUp(test.dir, func(path string) error {
			exist = append(exist, path)

			return nil
		})

		if !reflect.DeepEqual(exist, test.expected) {
			t.Error("filed to walk througn the test directory", test.dir, ".\r\n",
				"The list of exist child directories are\r\n", strings.Join(exist, ":"), "\r\n",
				", but expect directories are\r\n", strings.Join(test.expected, ":"))
		}
	}
}

func TestFilesRemoveAll(t *testing.T) {
	testRootDir := "data"
	testData := struct {
		path   []string
		check  []string
		files  []string
		remove string
	}{
		path: filepath.SplitList("data/test/path:data/test/other.path/пример"),
		check: filepath.SplitList("data:data/test"),
		files: []string{
			"some.test.txt",
			"пример.файла",
		},
		remove: "data/test/path",
	}

	// the list of check directories has extended with finish paths
	testData.check = append(testData.check, testData.path...)

	defer os.RemoveAll(testRootDir)

	// setup the test data
	for _, dirName := range testData.path {
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			t.Error("filed to create the test directory", dirName, ":", err)

			continue
		}

		tools.DirWalkUp(dirName, func(path string) error {
			var (
				fullPath string
			)

			for _, fileName := range testData.files {
				fullPath = filepath.Join(path, fileName)

				os.Create(fullPath)
			}

			return nil
		})
	}
	// check the test data
	walkTestFile := func(walkFunc func(string) bool) bool {
		for _, dirName := range testData.check {
			for _, fileName := range testData.files {
				if !walkFunc(filepath.Join(dirName, fileName)) {
					return false
				}
			}
		}

		return true
	}

	if !walkTestFile(func(fullPath string) bool {
		if info, err := os.Stat(fullPath); err != nil || !info.Mode().IsRegular() || info.IsDir() {
			t.Error("failed to check the test file structure. The file ", fullPath, " is not exist :", err)

			return false
		}

		return true
	}) {
		return
	}

	// execute test
	tools.FilesRemoveAll(testData.remove)

	// check the exist file structure
	walkTestFile(func(fullPath string) bool {
		if _, err := os.Stat(fullPath); err == nil {
			t.Error("failed to check the result file structure. The file ", fullPath, " is exist")

			return false
		}

		return true
	})
}