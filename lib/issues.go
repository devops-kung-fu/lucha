package lib

import (
	"bufio"
	"path/filepath"
)

func FindIssues(fs FileSystem, maxSeverity int) (violations []ScanFile, violationsDetected bool, err error) {
	files, err := BuildFileList(fs)
	if err != nil {
		return nil, false, err
	}
	for _, fl := range files {
		var scanFile ScanFile
		filename, _ := filepath.Abs(fl)
		if err != nil {
			return nil, false, err
		}
		file, err := fs.fs.Open(filename)
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
				// dir, _ := os.Getwd()
				// fn := strings.ReplaceAll(fl, dir, "")
				scanFile.Path = fl
				scanFile.Issues = append(scanFile.Issues, issues...)
				violationsDetected = true
			}
		}
		if violationsDetected {
			violations = append(violations, scanFile)
		}
	}
	return
}
