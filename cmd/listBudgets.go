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
	listBudgetsCmd.Flags().StringP("format", "f", "text", "Defines the rendering format")

}

func RenderBudgetsListText(budgets service.Budgets) {
	table := table.NewTable()
	table.SetHeader([]string{"CATEGORY", "VALUE"})
	var sum float64
	for k, v := range budgets {
		sum += v
		table.Append([]string{k, fmt.Sprintf("%.2f", v)})
	}
	table.AppendSeparator()
	table.Append([]string{"TOTAL", fmt.Sprintf("%0.2f", sum)})
	table.Render(os.Stdout) // Send output
}
