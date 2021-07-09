package lib

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

//LuchaDir returns the path where the default Lucha rules are stored
func LuchaDir() (path string, err error) {
	d, err := os.UserHomeDir()
	if err != nil {
		return
	}
	luchaDir := fmt.Sprintf("%s/%s", d, ".lucha")
	return luchaDir, nil
}

func Evaluate(line string, lineNumber int, maxSeverity int) (issues []Issue, err error) {
	for _, r := range Rules {
		if r.Severity >= int64(maxSeverity) {
			var issue Issue
			compiledRegex, err := regexp.Compile(r.Regex)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("%s has an invalid regex: %s", r.Code, r.Regex))
			}
			match := compiledRegex.Match([]byte(line))
			if match {
				issue = Issue{
					LineNumber: lineNumber,
					Rule:       r,
				}
				issues = append(issues, issue)
			}
		}
	}
	return
}

//RefreshRules pulls down the latest rules from https://github.com/devops-kung-fu/lucha
func RefreshRules(fs FileSystem, version string) (config Configuration, err error) {
	luchaDir, _ := LuchaDir()

	exists, _ := fs.Afero().DirExists(luchaDir)
	if !exists {
		err = fs.Afero().Mkdir(luchaDir, 0700)
		if err != nil {
			return
		}
	}
	_, err = DownloadFile(fs, luchaDir, "https://raw.githubusercontent.com/devops-kung-fu/lucha/main/lucha.yaml")
	if err != nil {
		return
	}
	return
}

//LoadRules loads the lucha.yaml rules file into memory
func LoadRules(fs FileSystem, version string, luchaFile string) (config Configuration, err error) {
	yamlFile, err := fs.Afero().ReadFile(luchaFile)
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
