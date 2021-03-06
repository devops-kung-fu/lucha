package cmd

import (
	"fmt"

	"github.com/devops-kung-fu/lucha/lib"
	"github.com/spf13/cobra"
)

var (
	refreshCmd = &cobra.Command{
		Use:   "refresh",
		Short: "Retrieves the latest rules for lucha.",
		Long:  "Retrieves the latest rules for lucha from https://github.com/devops-kung-fu/lucha.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Refreshing rules from DKFM...")
			fmt.Println()
			_, err := lib.RefreshRules(lib.NewOsFs(), version)
			if lib.IsErrorBool(err, "[ERROR]") {
				return
			}
		},
	}
)

func init() {
	rulesCmd.AddCommand(refreshCmd)
}
