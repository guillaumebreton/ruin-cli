package cmd

import (
	"os"
	"strconv"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/guillaumebreton/ruin/util"
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
			for _, budget := range args {
				println(budget, budgetValue)
				if budgetDelete {
					ledger.DeleteBudget(budget)
				} else {
					v, err := strconv.ParseFloat(budgetValue, 64)
					if err != nil {
						util.ExitOnError(err, "Fail to pase the value")
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
