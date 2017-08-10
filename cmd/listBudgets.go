package cmd

import (
	"os"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/spf13/cobra"
)

// listBCmd represents the list command
var listBudgetsCmd = &cobra.Command{
	Use:   "budgets",
	Short: "List all the budgets",
	Run: func(cmd *cobra.Command, args []string) {
		budgets := ledger.GetBudgets()
		RenderBudgetsListText(budgets)
	},
}

func init() {
	listCmd.AddCommand(listBudgetsCmd)

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
