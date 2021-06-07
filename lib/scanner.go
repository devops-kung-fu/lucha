package lib

import (
	"bufio"
	"path/filepath"
)

var (
	//Rules contains the loaded rules from lucha.yaml
	Rules []Rule
	//IgnoreFiles contains the names of files that shouldn't be processed from the .luchaignore file
	IgnoreFiles []string
)

func buildFileList(path string, recurse bool) (files []ScanFile, err error) {
	if recurse {
		files, err = ScanFilesRecursive(path)
	} else {
		files, err = ScanFiles(path)
	}
	if err != nil {
		return nil, err
	}
	return
}

func (f FileSystem) FindIssues(path string, recurse bool, maxSeverity int) (violations []ScanFile, violationsDetected bool, err error) {
	files, err := buildFileList(path, recurse)
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
				issues, err := Evaluate(line, lineNumber, maxSeverity)
				if err != nil {
					return nil, false, err
				}

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
