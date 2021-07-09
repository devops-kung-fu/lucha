package lib

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

//LoadIgnore loads content in from the .luchaignore file
func LoadIgnore(fs FileSystem) (ignores []string, err error) {
	filename := fmt.Sprintf("%s/.luchaignore", fs.SearchPath)
	exists, err := fs.Afero().Exists(filename)
	if !exists {
		return ignores, errors.New("no ignore file exists")
	}
	file, err := fs.fs.Open(filename)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		if !strings.HasPrefix(text, "#") {
			if fs.SearchPath != "./" {
				text = fmt.Sprintf("%s/%s", fs.SearchPath, text)
			}
			ignores = append(ignores, text)
		}
	}
	return
}
