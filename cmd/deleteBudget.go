package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
)

// deleteBudgetCmd represents the delete command
var deleteBudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Delete a budget",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Please to provide a category")
			os.Exit(1)

		}
		ledger.DeleteBudget(args[0])
		ledger.Save()
	},
}

func init() {
	deleteCmd.AddCommand(deleteBudgetCmd)
}
