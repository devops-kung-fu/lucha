package cmd

import (
	"fmt"

	"github.com/devops-kung-fu/lucha/lib"
	"github.com/spf13/cobra"
)

var (
	ignoreCmd = &cobra.Command{
		Use:   "ignore [filename]",
		Short: "Adds the specified file to .luchaignore.",
		Long:  "Adds the specified file to .luchaignore so they won't be checked for rules",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Added %s to .luchaignore\n", args[0])
			if lib.IsErrorBool(lib.NewOsFs().AppendIgnore(args[0]), "[ERROR]") {
				return
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(ignoreCmd)
}
