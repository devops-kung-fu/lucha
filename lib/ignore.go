package lib

import (
	"bufio"
	"fmt"
	"strings"
)

//Ignores a collection of patterns to ignore
var Ignores []string

//LoadIgnore loads content in from the .luchaignore file
func LoadIgnore(fs FileSystem, root string) (err error) {
	filename := fmt.Sprintf("%s/.luchaignore", root)
	exists, _ := fs.Afero().Exists(filename)
	if !exists {
		return nil
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
	var ignores []string

	for scanner.Scan() {
		text := scanner.Text()
		if !strings.HasPrefix(text, "#") {
			if fs.SearchPath != "./" {
				text = fmt.Sprintf("%s/%s", fs.SearchPath, text)
			}
			ignores = append(ignores, text)
		}
	}

	Ignores = ignores
	return
}
