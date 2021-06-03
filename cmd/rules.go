package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rulesCmd = &cobra.Command{
		Use:   "rules",
		Short: "Operations to manage rules.",
		Long:  "Operations to manage rules.",
	}
)

func init() {
	rootCmd.AddCommand(rulesCmd)
}
