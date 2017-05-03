package cmd

import (
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a transaction",
}

func init() {
	RootCmd.AddCommand(showCmd)
}
