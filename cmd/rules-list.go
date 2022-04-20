package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/devops-kung-fu/lucha/lib"
	"github.com/devops-kung-fu/lucha/util"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists currently loadable rules.",
		Long:  "Lists currently loadable rules",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := lib.LoadRules(lib.NewOsFs(), version, LuchaRulesFile)
			if err != nil {
				RulesFileNotFound()
			}
			fmt.Printf("Loading %v rules from %s", len(lib.Rules), LuchaRulesFile)
			fmt.Println()
			if util.IsErrorBool(err, "[ERROR]") {
				return
			}
			if util.IsErrorBool(displayRules(), "[ERROR]") {
				return
			}
			fmt.Println()
			fmt.Printf("Rules Version: %s\n", config.Version)
			fmt.Printf("# of Rules:    %v\n", len(lib.Rules))
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
