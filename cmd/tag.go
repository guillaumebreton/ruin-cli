package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/util"
	"github.com/m1ome/leven"
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
				count = autotag()
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

func autotag() int {
	count := 0
	f := service.NewFilter()
	txs := ledger.GetTransactions(f)

	nonCatTxs := []*service.Transaction{}
	catTxs := map[string]map[string]int{}
	// map a description with the tag + count
	for _, tx := range txs {
		if tx.Category != "" {
			v, ok := catTxs[tx.Description]
			if !ok {
				v = map[string]int{}
			}
			c, ok := v[tx.Category]
			if !ok {
				v[tx.Category] = 1
			} else {
				v[tx.Category] = c + 1
			}
			catTxs[tx.Description] = v

		} else {
			nonCatTxs = append(nonCatTxs, tx)
		}
	}

	// look for similarities
	for _, noCatTx := range nonCatTxs {
		var tags map[string]int
		var previousLength = 99999
		for k, v := range catTxs {
			l := leven.Distance(noCatTx.Description, k)
			if l < previousLength && l < 10 {
				tags = v
				previousLength = l
			}
		}
		if len(tags) != 0 {
			if noCatTx.Category = max(tags); noCatTx.Category != "" {
				count++
				ledger.UpdateTransaction(noCatTx.Number, noCatTx)
			}
		}
	}

	return count
}

func max(categories map[string]int) string {
	total := 0
	for _, count := range categories {
		total += count
	}
	catToApply := ""
	var score float32 = 0
	for cat, count := range categories {
		percentage := (float32(count) / float32(total)) * 100
		if percentage > score {
			score = percentage
			catToApply = cat
		}
	}
	if score > 33 && catToApply != "" && total != 2 {
		return catToApply
	}
	return ""
}

func init() {
	RootCmd.AddCommand(tagCmd)

	txCmd.Flags().BoolVarP(&autotagging, "auto-tagging", "a", true, "enable autotagging")
}
