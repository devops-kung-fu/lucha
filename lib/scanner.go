package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func FindIssues(path string, recurse bool) (violations []ScanFile, violationsDetected bool, err error) {
	var files []ScanFile
	if recurse {
		files, err = ScanFilesRecursive(path)
	} else {
		files, err = ScanFiles(path)
	}
	if err != nil {
		return nil, false, err
	}
	for _, f := range files {
		if IsTextFile(f) {
			filename, _ := filepath.Abs(f.Path)
			if err != nil {
				return nil, false, err
			}
			file, err := os.Open(filename)
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
					f.Issues = append(f.Issues, issues...)
					violationsDetected = true
				}
			}
			if violationsDetected {
				violations = append(violations, f)
			}
		}
	}
	return
}

func IsTextFile(file ScanFile) bool {

	buf, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s", file.Path, file.Info.Name()))
	size := 0
	for start := 0; start < len(buf); start += size {
		var r rune
		if r, size = utf8.DecodeRune(buf[start:]); r == utf8.RuneError {
			return false
		}
	}
	return true
}

func ScanFiles(path string) (files []ScanFile, err error) {
	tempFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range tempFiles {
		if !f.IsDir() && !Contains(IgnoreFiles, f.Name()) {
			files = append(files, ScanFile{
				Path: path,
				Info: f,
			})
		}
	}
	return
}

func ScanFilesRecursive(path string) (files []ScanFile, err error) {
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && !Contains(IgnoreFiles, info.Name()) {
				files = append(files, ScanFile{
					Path: path,
					Info: info,
				})
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	// err = filepath.Walk(path,
	// 	func(path string, info os.FileInfo, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if !info.IsDir() {

	// 		}
	// 		return nil
	// 	})
	return
}

func LineCounter(r io.Reader) (int, error) {

	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}
