package lib

import "regexp"

func Evaluate(line string, lineNumber int) (issues []Issue) {
	for _, r := range Rules {
		var issue Issue
		compiledRegex := regexp.MustCompile(r.Regex)

		match := compiledRegex.Match([]byte(line))
		if match {
			issue = Issue{
				LineNumber:  lineNumber,
				Description: r.Description,
				Severity:    3,
			}
			issues = append(issues, issue)
		}
	}
	return
}
