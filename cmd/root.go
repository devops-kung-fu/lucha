// Package cmd contains all of the commands that may be executed in the cli
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/devops-kung-fu/lucha/lib"
	"github.com/devops-kung-fu/lucha/util"
)

var (
	Afs            = &afero.Afero{Fs: afero.NewOsFs()}
	version        = "1.0.0"
	debug          bool
	Verbose        bool
	LuchaRulesFile string
	rootCmd        = &cobra.Command{
		Use:     "lucha",
		Short:   `"Scans for sensitive data in source code`,
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !debug {
				log.SetOutput(ioutil.Discard)
			}
			util.DoIf(Verbose, func() {
				fmt.Println()
				color.Style{color.FgWhite, color.OpBold}.Println("█   █ █ █▀▀ █ █ ▄▀█")
				color.Style{color.FgWhite, color.OpBold}.Println("█▄▄ █▄█ █▄▄ █▀█ █▀█")
				fmt.Println()
				fmt.Println("DKFM - DevOps Kung Fu Mafia")
				fmt.Println("https://github.com/devops-kung-fu/lucha")
				fmt.Printf("Version: %s\n", version)
				fmt.Println()
			})
		},
	}
)

// Execute creates the command tree and handles any error condition returned
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	luchaDir, _ := lib.LuchaDir()
	defaultLuchaFile := fmt.Sprintf("%s/lucha.yaml", luchaDir)
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "show debug output")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", true, "show verbose output")
	rootCmd.PersistentFlags().StringVar(&LuchaRulesFile, "rules-file", defaultLuchaFile, "Rules file to use when running lucha")
}

//RulesFileNotFound prints an error message if the lucha.yaml file isn't found and exits
func RulesFileNotFound() {
	util.PrintErr("[ERROR]", errors.New("NO RULES FOUND"))
	fmt.Println()
	util.PrintTabbed("The lucha rules file was not found. Ensure you have")
	util.PrintTabbed("run `lucha rules refresh` to install the rules locally, or")
	util.PrintTabbed("used the --rules-file flag to specify the rules to use.")
	os.Exit(1)
}
