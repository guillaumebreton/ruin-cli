package cmd

import (
	"fmt"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var endDate string
var startDate string
var listCategory string

// listTransactionCmd represents the listTransaction command
var listTransactionsCmd = &cobra.Command{
	Use:   "tx",
	Short: "List transactions",
	Run: func(cmd *cobra.Command, args []string) {
		f := service.NewFilter()
		if endDate != "" {
			t, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid end-date format")
				os.Exit(1)
			}
			f.EndDate = t
		}
		if startDate != "" {
			t, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid start-date format")
				os.Exit(1)
			}
			f.StartDate = t
		}

		if listCategory != "" {
			f.Category = listCategory
		}
		RenderTransactionListText(ledger.GetTransactions(f))
	},
}

func RenderTransactionListText(transactions []*service.Transaction) {
	table := table.NewTable()
	table.SetHeader([]string{"#", "DATE", "DESCRIPTION", "CATEGORY", "AMOUNT"})
	for _, v := range transactions {
		table.Append([]string{fmt.Sprintf("%d", v.Number), v.GetDate().Format("2006-01-02"), v.Description, v.Category, fmt.Sprintf("%.2f", v.Amount)})
	}
	table.Render(os.Stdout) // Send output
}

func init() {
	listCmd.AddCommand(listTransactionsCmd)

	listTransactionsCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "the start date")
	listTransactionsCmd.Flags().StringVarP(&endDate, "end-date", "e", "", "the end date")
	listTransactionsCmd.Flags().StringVarP(&listCategory, "category", "c", "", "category")

}
