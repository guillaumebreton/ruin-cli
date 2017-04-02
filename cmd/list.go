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
		budgets := service.GetBudgets()
		RenderListText(budgets)
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().StringP("format", "f", "text", "Defines the rendering format")

}

func RenderListText(budgets []service.Budget) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Value"})
	table.SetBorder(false)
	table.SetRowLine(false) // Enable row line
	table.SetHeaderLine(false)

	// Change table lines
	table.SetCenterSeparator("*")
	table.SetColumnSeparator("")
	table.SetRowSeparator("-")
	for _, v := range budgets {
		table.Append([]string{v.Category, fmt.Sprintf("%.2f", v.Value)})
	}
	table.Render() // Send output
}
