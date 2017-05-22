package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/util"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag a transaction",
	Run: func(cmd *cobra.Command, args []string) {
		// loop
		count := 0
		f := service.NewFilter()
		f.NoCategory = true
		tx := ledger.GetTransactions(f)
		for _, v := range tx {
			fmt.Printf("Enter a category for \"%s - %s - %s\" ? (skip/exit) \n", v.GetDate().Format("2006-01-02"), v.Description, fmt.Sprintf("%.2f", v.Amount))
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			util.ExitOnError(err, "Error : fail to read")
			switch input {
			case "skip\n", "\n":
				continue
			case "exit\n":
				goto exitLoop
				break
			default:
				v.Category = strings.TrimSuffix(input, "\n")
				ledger.UpdateTransaction(v.Number, v)
				count++
			}

		}
	exitLoop:
		if count > 0 {
			err := ledger.Save(ledgerFile)
			util.ExitOnError(err, "Fail to save ledger")
		}
		fmt.Printf("%d transaction(s) modified\n", count)
	},
}

func init() {
	RootCmd.AddCommand(tagCmd)
}
