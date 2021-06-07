package lib

import (
	"fmt"
	"os"
	"strings"
)

func PrintIf(f func(), condition bool) {
	if condition {
		f()
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func StartsWith(s []string, prefix string) bool {
	for _, v := range s {
		if strings.HasPrefix(v, prefix) {
			return true
		}
	}
	return false
}

func LuchaDir() (path string, err error) {

	d, err := os.UserHomeDir()
	if err != nil {
		return
	}
	luchaDir := fmt.Sprintf("%s/%s", d, ".lucha")
	return luchaDir, nil
}
