package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a budget",
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
