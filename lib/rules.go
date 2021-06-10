package lib

import (
	"errors"
	"fmt"
	"os"
	"regexp"
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
