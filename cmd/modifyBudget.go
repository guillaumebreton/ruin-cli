package cmd

import (
	"github.com/guillaumebreton/ruin/util"
	"github.com/spf13/cobra"
	"strconv"
)

// setBudgetCmd represents the set command
var modifyBudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Modify a budget",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			util.Exit("Need to provide a category and a value")
		}
		key := args[0]
		value := args[1]
		v, _ := strconv.ParseFloat(value, 64)
		ledger.SetBudget(key, v)
		err := ledger.Save(ledgerFile)
		util.ExitOnError(err, "Fail to save ledger file")
	},
}

func init() {
	modifyCmd.AddCommand(modifyBudgetCmd)
}
