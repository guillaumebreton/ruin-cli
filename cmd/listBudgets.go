package cmd

import (
	"fmt"
	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listBudgetsCmd = &cobra.Command{
	Use:   "budgets",
	Short: "List all the budgets",
	Long:  `Define`,
	Run: func(cmd *cobra.Command, args []string) {
		l, err := service.LoadLedger()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err: %v\n", err)
			os.Exit(1)
		}
		budgets := l.GetBudgets()
		RenderBudgetsListText(budgets)
	},
}

func init() {
	listCmd.AddCommand(listBudgetsCmd)

}

func RenderBudgetsListText(budgets map[string]service.Budget) {
	table := table.NewTable()
	table.SetHeader([]string{"CATEGORY", "FROM", "TO"})
	for k, v := range budgets {
		table.Append([]string{k, fmt.Sprintf("%s", v.StartDate.ToString()), fmt.Sprintf("%s", v.EndDate.ToString())})
	}
	table.Render(os.Stdout)
}
