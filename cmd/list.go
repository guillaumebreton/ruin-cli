package cmd

import (
	"fmt"
	"github.com/guillaumebreton/gobud/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the budgets",
	Long:  `Define`,
	Run: func(cmd *cobra.Command, args []string) {
		c, _ := service.LoadConfig()
		budgets := c.GetBudgets()
		RenderListText(budgets)
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("format", "f", "text", "Defines the rendering format")

}

func RenderListText(budgets service.Budgets) {
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
