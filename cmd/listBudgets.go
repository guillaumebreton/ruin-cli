package cmd

import (
	"fmt"
	"github.com/guillaumebreton/ruin/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listBudgetsCmd = &cobra.Command{
	Use:   "list",
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
	budgetCmd.AddCommand(listBudgetsCmd)
	listBudgetsCmd.Flags().StringP("format", "f", "text", "Defines the rendering format")

}

func RenderBudgetsListText(budgets service.Budgets) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Value"})
	table.SetBorder(false)
	table.SetRowLine(false) // Enable row line
	table.SetHeaderLine(false)

	// Change table lines
	table.SetCenterSeparator("*")
	table.SetColumnSeparator("")
	table.SetRowSeparator("-")
	for k, v := range budgets {
		table.Append([]string{k, fmt.Sprintf("%.2f", v)})
	}
	table.Render() // Send output
}
