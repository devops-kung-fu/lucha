package lib

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var (
	Rules       []Rule
	IgnoreFiles []string
)

type FileSystem struct {
	fs afero.Fs
}

func NewOsFs() FileSystem {
	var d FileSystem
	d.fs = afero.NewOsFs()
	return d
}

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
func (f FileSystem) LoadIgnore() (err error) {

	filename, _ := filepath.Abs(".luchaignore")
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
		ignores = append(ignores, scanner.Text())
	}

	IgnoreFiles = ignores
	return
}

func (f FileSystem) LoadRules(version string) (config Configuration, err error) {
	err = f.LoadIgnore()
	if err != nil {
		return
	}
	filename, _ := filepath.Abs("lucha.yaml")
	yamlFile, err := f.Afero().ReadFile(filename)

	if err != nil {
		return
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return
	}

	err = checkVersion(config, version)

	Rules = config.Lucha.Rules

	return
}

func checkVersion(config Configuration, version string) (err error) {
	if config.Version == "" {
		err = errors.New("no version value found in lucha.yaml")
		return
	}
	if version == "" {
		err = errors.New("version should not be empty")
		return
	}
	ver := strings.Split(config.Version, ".")
	verMatch := strings.Split(version, ".")
	if fmt.Sprintf("%v.%v", ver[0], ver[1]) != fmt.Sprintf("%v.%v", verMatch[0], verMatch[1]) {
		err = fmt.Errorf("version mismatch: update your lucha rules file (lucha.yaml)")
	}
	return
}

type Configuration struct {
	Version string `json:"version"`
	Lucha   Lucha  `json:"lucha"`
}

type Lucha struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Message     string `json:"message"`
	Attribution string `json:"attribution,omitempty"`
	Regex       string `json:"regex"`
	Severity    int64  `json:"severity,omitempty"`
}

type ScanFile struct {
	Path   string
	Info   fs.FileInfo
	Issues []Issue
}

//Returns the number of Issues in a scanned file
func (s *ScanFile) IssueCount() int {
	return len(s.Issues)
}

type Issue struct {
	LineNumber int
	Rule       Rule
}
