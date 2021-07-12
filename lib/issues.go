package lib

import (
	"bufio"
	"path/filepath"
	"strings"
)

//FindIssues scans the provided filesystem for issues
func FindIssues(fs FileSystem, minSeverity int) (violations []ScanFile, violationsDetected bool, err error) {
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
		if isUTF8(fs, file) { //Only scan UTF8 files

			lineNumber := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				lineNumber++
				issues, err := Evaluate(line, lineNumber, minSeverity)
				if err != nil {
					return nil, false, err
				}

				if len(issues) > 0 {
					scanFile.Path = strings.ReplaceAll(fl, fs.AbsoluteSearchPath(), strings.TrimSuffix(fs.SearchPath, "/"))
					scanFile.Issues = append(scanFile.Issues, issues...)
					violationsDetected = true
				}
			}
			if violationsDetected {
				violations = append(violations, scanFile)
			}
		}

		// this could go into a verbose or trace flag
		// else {
		// 	fmt.Println("Ignoring ", file.Name())
		// }
	}
	return
}
