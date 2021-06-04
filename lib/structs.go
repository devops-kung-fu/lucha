package lib

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

type Configuration struct {
	Version string `json:"version"`
	Lucha   Lucha  `json:"lucha"`
}

func (c *Configuration) checkVersion(version string) (err error) {
	if c.Version == "" {
		err = errors.New("no version value found in lucha.yaml")
		return
	}
	if version == "" {
		err = errors.New("version should not be empty")
		return
	}
	ver := strings.Split(c.Version, ".")
	verMatch := strings.Split(version, ".")
	if fmt.Sprintf("%v.%v", ver[0], ver[1]) != fmt.Sprintf("%v.%v", verMatch[0], verMatch[1]) {
		err = fmt.Errorf("version mismatch: update your lucha rules file (lucha.yaml)")
	}
	return
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
