package lib

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var (
	Rules []Rule
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

func (f FileSystem) LoadRules(version string) (config Configuration, err error) {
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
	Description string `json:"description"`
	Regex       string `json:"regex"`
	Severity    int64  `json:"severity"`
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
	LineNumber  int
	Description string
	Severity    int
}
