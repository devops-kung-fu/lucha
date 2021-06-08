package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var (
	//Rules contains the loaded rules from lucha.yaml
	Rules []Rule
	//IgnoreFiles contains the names of files that shouldn't be processed from the .luchaignore file
	IgnoreFiles []string
)

//FileSystem encapsulates the Afero fs Filesystem
type FileSystem struct {
	fs afero.Fs
}

//NewOsFs returns a new local os file system
func NewOsFs() FileSystem {
	var d FileSystem
	d.fs = afero.NewOsFs()
	return d
}

//Afero returns the Afero system
func (f FileSystem) Afero() (afs *afero.Afero) {
	afs = &afero.Afero{Fs: f.fs}
	return
}

//AppendIgnore appends the provided filename to the .luchaignore file
func (f FileSystem) AppendIgnore(filename string) (err error) {
	fn, _ := filepath.Abs(".luchaignore")
	contains, _ := f.Afero().FileContainsBytes(fn, []byte(fmt.Sprintf("%s\n", filename)))
	if !contains {
		file, err := f.fs.OpenFile(fn,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer func() {
			err = file.Close()
		}()

		if _, err := file.WriteString(fmt.Sprintf("%s\n", filename)); err != nil {
			return err
		}
	}
	return
}

//LoadIgnore loads content in from the .luchaignore file
func (f FileSystem) LoadIgnore() (err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	filename := fmt.Sprintf("%s/.luchaignore", pwd)
	exists, _ := f.Afero().Exists(filename)
	if !exists {
		return nil
	}
	file, err := f.fs.Open(filename)

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
			ignores = append(ignores, text)
		}
	}

	IgnoreFiles = ignores
	return
}

//LuchaRulesFile when passed a filename returns the absolute path
func (f FileSystem) LuchaRulesFile(file string) (luchaFile string, err error) {
	if filepath.IsAbs(file) {
		luchaFile = file
		return
	} else {
		luchaFile, err = filepath.Abs(file)
		return
	}
}

//LoadRules loads the lucha.yaml rules file into memory
func (f FileSystem) LoadRules(version string, file string) (config Configuration, err error) {
	filename, err := f.LuchaRulesFile(file)
	if err != nil {
		return
	}
	err = f.LoadIgnore()
	if err != nil {
		return
	}

	yamlFile, err := f.Afero().ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return
	}

	err = config.checkVersion(version)

	Rules = config.Lucha.Rules

	return
}

//RefreshRules pulls down the latest rules from https://github.com/devops-kung-fu/lucha
func (f FileSystem) RefreshRules(version string) (config Configuration, err error) {
	luchaDir, _ := LuchaDir()

	exists, _ := f.Afero().DirExists(luchaDir)
	if !exists {
		err = f.Afero().Mkdir(luchaDir, 0700)
		if err != nil {
			return
		}
	}
	_, err = DownloadURL("https://raw.githubusercontent.com/devops-kung-fu/lucha/main/lucha.yaml", luchaDir)
	if err != nil {
		return
	}
	return
}

//IsTextFile examines a file and returns true if the file is UTF-8
func (f FileSystem) IsTextFile(file ScanFile) bool {
	buf, _ := f.Afero().ReadFile(fmt.Sprintf("%s/%s", file.Path, file.Info.Name()))
	size := 0
	for start := 0; start < len(buf); start += size {
		var r rune
		if r, size = utf8.DecodeRune(buf[start:]); r == utf8.RuneError {
			return false
		}
	}
	return true
}

//ScanFiles grabs a list of files from the provided directory for scanning
func (f FileSystem) ScanFiles(path string) (files []ScanFile, err error) {
	// absPath, err := filepath.Abs(path)
	// if err != nil {
	// 	return nil, err
	// }
	tempFiles, err := f.Afero().ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range tempFiles {
		// if !f.IsDir() && !Contains(IgnoreFiles, f.Name()) {
		// 	if StartsWith(IgnoreFiles, path) {
		// 		break
		// 	}
		if !f.IsDir() {
			files = append(files, ScanFile{
				Path: path,
				Info: f,
			})
		}
		// }
	}
	return
}

//ScanFilesRecursive grabs a list of all files recursively for scanning
func (f FileSystem) ScanFilesRecursive(path string) (files []ScanFile, err error) {
	err = filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && !Contains(IgnoreFiles, info.Name()) {
				if StartsWith(IgnoreFiles, path) {
					return nil
				}
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
	return
}

// BuildFileList gathers the files to scan
func (f FileSystem) BuildFileList(path string, recurse bool) (files []ScanFile, err error) {
	if recurse {
		files, err = f.ScanFilesRecursive(path)
	} else {
		files, err = f.ScanFiles(path)
	}
	if err != nil {
		return nil, err
	}
	return
}

func (f FileSystem) FindIssues(path string, recurse bool, maxSeverity int) (violations []ScanFile, violationsDetected bool, err error) {
	files, err := f.BuildFileList(path, recurse)
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
