package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/devops-kung-fu/lucha/lib"
	"github.com/devops-kung-fu/lucha/util"
)

var (
	includeGit  bool
	recursive   bool
	minSeverity int
	scanCmd     = &cobra.Command{
		Use:     "scan",
		Short:   "Scans all files in the recursively and scans them for sensitive data.",
		Long:    "Scans all files in the recursively and scans them for sensitive data.",
		Example: "lucha scan --recursive .",
		Run: func(cmd *cobra.Command, args []string) {
			fs := lib.NewOsFs()

			if len(args) == 0 {
				color.Style{color.FgRed, color.OpBold}.Println("Please provide the path to the repository.")
				fmt.Println()
				_ = cmd.Usage()
			} else if len(args) > 1 {
				color.Style{color.FgRed, color.OpBold}.Println("Only one path is allowed.")
				fmt.Println()
				_ = cmd.Usage()
			} else {
				fs.SearchPath = args[0]
				if !strings.HasSuffix(fs.SearchPath, "/") {
					fs.SearchPath = fmt.Sprintf("%s/", fs.SearchPath)
				}

				fs.Recursive = recursive
				fs.IncludeGit = includeGit

				err := initScan(fs)

				if util.IsErrorBool(err, "[ERROR]") {
					os.Exit(1)
				}

				s := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
				s.Start()
				fmt.Printf("Scanning files in %s\n\n", fs.SearchPath)
				files, issuesDetected, err := lib.FindIssues(fs, minSeverity)
				s.Stop()
				util.IfErrorLog(err, "[ERROR]")
				if issuesDetected {
					color.Style{color.FgRed.Darken()}.Println("ISSUES DETECTED!")
					fmt.Println()
					for _, f := range files {
						if len(f.Issues) > 0 {
							fmt.Println(f.Path)
							for _, i := range f.Issues {
								fmt.Print("  ")
								printSeverityIndicator(int(i.Rule.Severity))
								fmt.Printf(" %s:%v:1, %s\n", f.Path, i.LineNumber, i.Rule.Message)
							}
						}
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
	scanCmd.PersistentFlags().IntVar(&minSeverity, "min-severity", 0, "Only report on severities higher than this value")
	scanCmd.PersistentFlags().BoolVarP(&includeGit, "git", "g", false, "If true, lucha not ignore the .git directory")
}

func initScan(fs lib.FileSystem) (err error) {

	// err = lib.LoadIgnore(fs, fs.SearchPath)
	// if err != nil {
	// 	return
	// }

	_, err = lib.LoadRules(fs, version, LuchaRulesFile)
	if err != nil {
		RulesFileNotFound()
	}
	fmt.Printf("%v Rules Loaded\n\n", len(lib.Rules))
	util.IfErrorLog(err, "[ERROR]")
	return
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
