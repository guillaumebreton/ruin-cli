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

var autotagging bool

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Interactive command to tag process",
	Run: func(cmd *cobra.Command, args []string) {
		total := 0
		count := 1
		if autotagging {
			for count > 0 {
				count = ledger.Autotag()
				total += count
			}
		}
		total += interactive()
		if total > 0 {
			err := ledger.Save(ledgerFile)
			util.ExitOnError(err, "Fail to save ledger")
		}

		println(total, " transactions have been tagged")
	},
}

func interactive() int {

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
			return count
		default:
			v.Category = strings.TrimSuffix(input, "\n")
			ledger.UpdateTransaction(v.Number, v)
			count++
		}

	}
	if ledger.Dirty {
		ledger.Save(ledgerFile)
	}
	return count
}

func init() {
	RootCmd.AddCommand(tagCmd)

	txCmd.Flags().BoolVarP(&autotagging, "auto-tagging", "a", true, "enable autotagging")
}
