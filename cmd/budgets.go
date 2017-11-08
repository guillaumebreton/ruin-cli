package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/guillaumebreton/ruin-cli/service"
	"github.com/guillaumebreton/ruin-cli/table"
	"github.com/guillaumebreton/ruin-cli/util"
	"github.com/spf13/cobra"
)

// budgetsCmd represents the budgets command
var (
	budgetValue  string
	budgetDelete bool
)
var budgetsCmd = &cobra.Command{
	Use:   "budgets",
	Short: "Budgets management command",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			budgets := ledger.GetBudgets()
			RenderBudgetsListText(budgets)
		} else {
			if !budgetDelete && budgetValue == "" {
				util.Exit("Please use -v with a value or -d")
			}
			for _, budget := range args {
				if budgetDelete {
					err := ledger.DeleteBudget(budget)
					util.ExitOnError(err, "Fail to delete budget")
					fmt.Println("Budget deleted")
				} else {
					var v float64
					var err error
					if budgetValue == "" {
						v = 0
					} else {
						v, err = strconv.ParseFloat(budgetValue, 64)
					}
					if err != nil {
						util.ExitOnError(err, "Fail to parese the value")
					}
					err = ledger.SetBudget(budget, v)
					if err != nil {
						util.ExitOnError(err, "Fail to set budget value")
					}
				}
			}
			ledger.Save(ledgerFile)
		}
	},
}

func init() {
	RootCmd.AddCommand(budgetsCmd)
	budgetsCmd.Flags().StringVarP(&budgetValue, "value", "v", "", "Set the threshold value of the budget")
	budgetsCmd.Flags().BoolVarP(&budgetDelete, "delete", "d", false, "delete the budget")
}

func RenderBudgetsListText(budgets service.Budgets) {
	table := table.NewTable()
	table.SetHeader("CATEGORY", "VALUE")
	var sum float64
	for k, v := range budgets {
		sum += v
		table.Append(k, v)
	}
	table.AppendSeparator()
	table.Append("TOTAL", sum)
	table.Render(os.Stdout)
}
