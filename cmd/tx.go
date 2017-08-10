package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/guillaumebreton/ruin/util"
	"github.com/spf13/cobra"
)

// ruin tx
// ruin tx 24 -c test
// ruin tx 24,25 -c test -d 12

var (
	txCategory  string
	txDate      string
	txEndDate   string
	txStartDate string
)

// txCmd represents the tx command
var txCmd = &cobra.Command{
	Use:   "tx [comma separated tx ids]",
	Short: "General tx command to list or edit transactions",
	Long: ` General tx command to list or edit transactions.
	usage:
	- ruin tx -s 2016-02-03 # list transactions after the given start date
	- ruin tx 23,34 -c books # set the category to books for 23 and 34
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			f := service.NewFilter()
			if txEndDate != "" {
				t, err := time.Parse("2006-01-02", txEndDate)
				util.ExitOnError(err, "Invalid end-date format")
				f.EndDate = t
			}
			if txStartDate != "" {
				t, err := time.Parse("2006-01-02", txStartDate)
				util.ExitOnError(err, "Invalid start-date format")
				f.StartDate = t
			}

			if txCategory != "" {
				f.Category = txCategory
			}
			// reindex if the ledger is dirty
			if ledger.Dirty {
				ledger.Reindex()
				err := ledger.Save(ledgerFile)
				util.ExitOnError(err, "Fail to save file")
			}
			RenderTransactionListText(ledger.GetTransactions(f))
		} else {

			arr := strings.Split(args[0], ",")

			var err error
			if txDate == "" && txCategory == "" {
				for _, v := range arr {
					id, err := strconv.Atoi(v)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%s not found\n", v)
					} else {
						tx, err := ledger.GetTransactionByNumber(id)
						util.ExitOnError(err, "Transaction not found")
						RenderShowTransaction(tx)
					}
				}
			} else {
				// if -c or -d is set than modify the date
				var d time.Time
				if txDate != "" {
					d, err = time.Parse(service.ShortFormat, txDate)
					util.ExitOnError(err, "Invalid date")
				}
				for _, v := range arr {
					id, err := strconv.Atoi(v)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%s not found\n", v)
					} else {
						tx, err := ledger.GetTransactionByNumber(id)
						if err != nil {
							fmt.Fprintf(os.Stderr, "%d not found\n", id)
						} else {
							if txCategory != "" {
								tx.Category = txCategory
							}
							if txDate != "" {
								tx.UserDate = d
							}
							ledger.UpdateTransaction(id, tx)
							fmt.Printf("Transaction %d updated\n", id)
						}
					}
				}
				err = ledger.Save(ledgerFile)
				util.ExitOnError(err, "Fail to save file")
			}

		}
	},
}

func RenderShowTransaction(tx *service.Transaction) {
	fmt.Printf("\nField       | Values \n")
	fmt.Printf("------------ ----------------------------\n")
	fmt.Printf("Number      : %d\n", tx.Number)
	fmt.Printf("ID          : %s\n", tx.ID)
	fmt.Printf("Date        : %s\n", tx.Date.Format(service.ShortFormat))
	fmt.Printf("User date   : %s\n", tx.UserDate.Format(service.ShortFormat))
	fmt.Printf("Amount      : %0.2f\n", tx.Amount)
	fmt.Printf("Category    : %s\n", tx.Category)
	fmt.Printf("Balance     : %0.2f\n\n", tx.Balance)
}
func RenderTransactionListText(transactions []*service.Transaction) {
	table := table.NewTable()
	table.SetHeader("#", "DATE", "DESCRIPTION", "CATEGORY", "AMOUNT")
	for _, v := range transactions {
		table.Append([]string{fmt.Sprintf("%d", v.Number), v.GetDate().Format("2006-01-02"), v.Description, v.Category, fmt.Sprintf("%.2f", v.Amount)})
	}
	table.Render(os.Stdout) // Send output
}

func init() {
	RootCmd.AddCommand(txCmd)
	txCmd.Flags().StringVarP(&txCategory, "category", "c", "", "Set the transaction category")
	txCmd.Flags().StringVarP(&txDate, "date", "d", "", "Set the transaction date (eg. 2017-02-03)")
	txCmd.Flags().StringVarP(&txStartDate, "start-date", "s", "", "Filter the list with the start date (eg. 2017-02-03)")
	txCmd.Flags().StringVarP(&txEndDate, "end-date", "e", "", "Filter the list with the end date (eg. 2017-02-03)")
}
