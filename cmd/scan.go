package cmd

import (
	"fmt"
	"os"

	"github.com/devops-kung-fu/lucha/lib"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	recursive bool
	scanCmd   = &cobra.Command{
		Use: "scan [path]",

		Short: "Scans all files in the recursively and scans them for sensitive data.",
		Long:  "Scans all files in the recursively and scans them for sensitive data.",
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			_, err := lib.NewOsFs().LoadRules(version)
			fmt.Printf("%v Rules Loaded\n\n", len(lib.Rules))
			lib.IfErrorLog(err, "[ERROR]")
			fmt.Printf("Scanning files in %s\n\n", path)
			files, issuesDetected, err := lib.FindIssues(path, recursive)
			lib.IfErrorLog(err, "[ERROR]")
			if issuesDetected {
				color.Style{color.FgRed}.Println("ISSUES DETECTED!")
				fmt.Println()
				for _, f := range files {
					if f.IssueCount() > 0 {
						fmt.Println(f.Path)
						for _, i := range f.Issues {
							fmt.Printf("     %s:%v:1, %s\n", f.Path, i.LineNumber, i.Description)
						}
					}
				}
				os.Exit(1)
			} else {
				color.Style{color.FgLightGreen}.Println("No Issues Detected")
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "If true, lucha will recurse subdirectories")
}
