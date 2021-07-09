package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_toAbsolutePath(t *testing.T) {

	searchDir, _ := filepath.Abs("../../hookz")
	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range fileList {
		fmt.Println(file)
	}
}
