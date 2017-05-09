package cmd

import (
	"fmt"
	"github.com/guillaumebreton/ruin/service"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	modifyCategory string
	modifyDate     string
)

// modifyTransactionCmd represents the modify command
var modifyTransactionCmd = &cobra.Command{
	Use:   "tx",
	Short: "Modify a transaction",
	Run: func(cmd *cobra.Command, args []string) {
		var d time.Time
		var err error
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Please provide a list of transaction ids")
			os.Exit(1)
		}
		if modifyDate != "" {
			d, err = time.Parse(service.ShortFormat, modifyDate)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid date : %s", modifyDate)
				os.Exit(1)
			}

		}
		arr := strings.Split(args[0], ",")
		for _, v := range arr {
			id, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s not found", v)
			} else {
				tx, err := ledger.GetTransactionByNumber(id)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%d not found", id)
				} else {
					if modifyCategory != "" {
						tx.Category = modifyCategory
					}
					if modifyDate != "" {
						println("Set date")
						tx.UserDate = d
					}
					ledger.UpdateTransaction(id, tx)
				}
			}
		}
		fmt.Printf("Transaction %s updated\n", args[0])
		ledger.Save(ledgerFile)
	},
}

func init() {
	modifyCmd.AddCommand(modifyTransactionCmd)
	modifyTransactionCmd.Flags().StringVarP(&modifyCategory, "category", "c", "", "Set the transaction category")
	modifyTransactionCmd.Flags().StringVarP(&modifyDate, "date", "d", "", "Set the transaction date (eg. 2017-02-03)")
}
