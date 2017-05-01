package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
)

// deleteBudgetCmd represents the delete command
var deleteBudgetCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a budget on a category",
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
	budgetCmd.AddCommand(deleteBudgetCmd)
}
