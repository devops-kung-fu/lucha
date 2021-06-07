// Package cmd contains all of the commands that may be executed in the cli
package cmd

import (
	"fmt"
	"os"

	"github.com/devops-kung-fu/lucha/lib"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	version   = "0.0.1"
	verbose   bool
	RulesFile string
	rootCmd   = &cobra.Command{
		Use:     "lucha scan",
		Short:   `"Scans for sensitive data in source code`,
		Version: version,
	}
)

// Execute creates the command tree and handles any error condition returned
func Execute() {
	cobra.OnInitialize(func() {
		// var fs = afero.NewOsFs()
		// afs := &afero.Afero{Fs: fs}
		// b, err := afs.DirExists(".git")
		// lib.IfErrorLog(err, "[ERROR]")
		// if !b {
		// 	e := errors.New("*** must be run in a local .git repository")
		// 	lib.IfErrorLog(e, "ERROR")
		// 	os.Exit(1)
		// }
	})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	color.Style{color.FgWhite, color.OpBold}.Println("lucha")
	fmt.Println("https://github.com/devops-kung-fu/lucha")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("")

	luchaDir, _ := lib.LuchaDir()
	defaultLuchaFile := fmt.Sprintf("%s/lucha.yaml", luchaDir)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Extended output as lucha executes.")
	rootCmd.PersistentFlags().StringVar(&RulesFile, "rules-file", defaultLuchaFile, "Rules file to use when running lucha")
}
