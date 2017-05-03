package cmd

import (
	"fmt"
	"github.com/guillaumebreton/ruin/service"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// showTxCmd represents the showTx command
var showTxCmd = &cobra.Command{
	Use:   "tx",
	Short: "Show a transaction",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Please provide a transaction id")
			os.Exit(1)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s not found", args[0])
		} else {
			tx, err := ledger.GetTransaction(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%d not found", id)
				os.Exit(1)
			}
			RenderShowTransaction(tx)
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

func init() {
	showCmd.AddCommand(showTxCmd)
}
