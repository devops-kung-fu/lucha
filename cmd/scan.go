package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/devops-kung-fu/lucha/lib"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	recursive   bool
	maxSeverity int
	scanCmd     = &cobra.Command{
		Use:     "scan",
		Short:   "Scans all files in the recursively and scans them for sensitive data.",
		Long:    "Scans all files in the recursively and scans them for sensitive data.",
		Example: "lucha scan --recursive .",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				color.Style{color.FgRed, color.OpBold}.Println("Please provide the path to the repository.")
				fmt.Println()
				_ = cmd.Usage()
			} else if len(args) > 1 {
				color.Style{color.FgRed, color.OpBold}.Println("Only one path is allowed.")
				fmt.Println()
				_ = cmd.Usage()
			} else {
				path := args[0]
				fs := lib.NewOsFs()
				_, err := fs.LoadRules(version, RulesFile)
				if err != nil {
					RulesFileNotFound()
				}
				fmt.Printf("%v Rules Loaded\n\n", len(lib.Rules))
				lib.IfErrorLog(err, "[ERROR]")
				s := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
				s.Start()
				fmt.Printf("Scanning files in %s\n\n", path)
				files, issuesDetected, err := fs.FindIssues(path, recursive, maxSeverity)
				s.Stop()
				lib.IfErrorLog(err, "[ERROR]")
				if issuesDetected {
					color.Style{color.FgRed.Darken()}.Println("ISSUES DETECTED!")
					fmt.Println()
					for _, f := range files {
						if len(f.Issues) > 0 {
							fmt.Println(f.Path)
							for _, i := range f.Issues {
								fmt.Print("    ")
								printSeverityIndicator(int(i.Rule.Severity))
								fmt.Printf(" %s:%v:1, %s\n", f.Path, i.LineNumber, i.Rule.Message)
							}
						}
					}
					if NoFail {
						os.Exit(0)
					}
					os.Exit(1)
				} else {
					color.Style{color.FgLightGreen}.Println("No Issues Detected")
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "If true, lucha will recurse subdirectories")
	scanCmd.PersistentFlags().IntVar(&maxSeverity, "max-severity", 0, "Only report on severities higher than this value")
}

func printSeverityIndicator(severity int) {
	switch severity {
	case 0:
		color.Style{color.FgBlue}.Print("■")
	case 1:
		color.Style{color.FgCyan}.Print("■")
	case 2:
		color.Style{color.FgYellow}.Print("■")
	case 3:
		color.Style{color.FgMagenta}.Print("■")
	case 4:
		color.Style{color.FgRed}.Print("■")
	}

}
