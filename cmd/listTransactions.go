package cmd

import (
	"fmt"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/guillaumebreton/ruin/util"
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
			util.ExitOnError(err, "Invalid end-date format")
			f.EndDate = t
		}
		if startDate != "" {
			t, err := time.Parse("2006-01-02", startDate)
			util.ExitOnError(err, "Invalid end-date format")
			f.StartDate = t
		}

		if listCategory != "" {
			f.Category = listCategory
		}
		// reindex if the ledger is dirty
		if ledger.Dirty {
			ledger.Reindex()
			err := ledger.Save(ledgerFile)
			util.ExitOnError(err, "Fail to save file")
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
