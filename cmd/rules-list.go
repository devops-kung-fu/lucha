package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/devops-kung-fu/lucha/lib"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists currently loadable rules.",
		Long:  "Lists currently loadable rules",
		Run: func(cmd *cobra.Command, args []string) {
			fs := lib.NewOsFs()
			_, err := fs.LoadRules(version, RulesFile)
			if err != nil {
				RulesFileNotFound()
			}
			luchaRulesFile, _ := fs.LuchaRulesFile(RulesFile)
			fmt.Printf("Loading %v rules from %s", len(lib.Rules), luchaRulesFile)
			fmt.Println()
			if lib.IsErrorBool(err, "[ERROR]") {
				return
			}
			if lib.IsErrorBool(displayRules(), "[ERROR]") {
				return
			}
		},
	}
)

func init() {
	rulesCmd.AddCommand(listCmd)
}

func displayRules() (err error) {
	tmpl := genTemplate()
	err = tmpl.Execute(os.Stdout, lib.Rules)
	if err != nil {
		return err
	}
	return
}

func genTemplate() (t *template.Template) {

	content := `{{range .}}
Code: {{.Code}}
Name: {{.Name}}
Description: {{.Description}}
Attribution: {{.Attribution}}
Severity: {{.Severity}}
{{end}}`

	return template.Must(template.New("rule").Parse(content))
}
