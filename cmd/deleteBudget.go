package cmd

import (
	"github.com/guillaumebreton/ruin/util"
	"github.com/spf13/cobra"
)

// deleteBudgetCmd represents the delete command
var deleteBudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Delete a budget",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			util.Exit("Please to provide a category")
		}
		ledger.DeleteBudget(args[0])
		err := ledger.Save(ledgerFile)
		util.ExitOnError(err, "Fail to save ledger")
	},
}

func init() {
	deleteCmd.AddCommand(deleteBudgetCmd)
}
