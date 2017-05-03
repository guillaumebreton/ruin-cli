package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// setBudgetCmd represents the set command
var modifyBudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Modify a budget",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Need to provide a category and a value")
			os.Exit(1)
		}
		key := args[0]
		value := args[1]
		v, _ := strconv.ParseFloat(value, 64)
		ledger.SetBudget(key, v)
		ledger.Save(ledgerFile)
	},
}

func init() {
	modifyCmd.AddCommand(modifyBudgetCmd)
	modifyBudgetCmd.Flags().StringVarP(&modifyCategory, "category", "c", "", "Set the transaction category")
}
