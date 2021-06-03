package lib

import (
	"bufio"
	"path/filepath"
)

var (
	Rules       []Rule
	IgnoreFiles []string
)

func (f FileSystem) FindIssues(path string, recurse bool) (violations []ScanFile, violationsDetected bool, err error) {
	var files []ScanFile
	if recurse {
		files, err = ScanFilesRecursive(path)
	} else {
		files, err = ScanFiles(path)
	}
	if err != nil {
		return nil, false, err
	}
	for _, fl := range files {
		if f.IsTextFile(fl) {
			filename, _ := filepath.Abs(fl.Path)
			if err != nil {
				return nil, false, err
			}
			file, err := f.fs.Open(filename)
			defer func() {
				err = file.Close()
			}()
			if err != nil {
				return nil, false, err
			}
			lineNumber := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				lineNumber++
				issues := Evaluate(line, lineNumber)

				if len(issues) > 0 {
					fl.Issues = append(fl.Issues, issues...)
					violationsDetected = true
				}
			}
			if violationsDetected {
				violations = append(violations, fl)
			}
		}
	}
	return
}
