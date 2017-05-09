package cmd

import (
	"fmt"
	"github.com/guillaumebreton/ruin/ofx"
	"github.com/guillaumebreton/ruin/util"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import an ofx file into the ledger",
	Run: func(cmd *cobra.Command, args []string) {

		o, err := ofx.Parse(args[0])
		util.ExitOnError(err, "Fail to parse OFX file")
		count := 0
		for _, tx := range o.Transactions {
			a := tx.GetAmount()
			if ledger.Add(tx.ID, tx.GetDate(), tx.TxType, tx.Description, a) {
				count++
			}
		}
		ledger.Balance = o.Balance
		err = ledger.Save(ledgerFile)
		util.ExitOnError(err, "Fail to save ledger")
		fmt.Printf("%d transaction(s) added\n", count)
	},
}

func init() {
	RootCmd.AddCommand(importCmd)
}
