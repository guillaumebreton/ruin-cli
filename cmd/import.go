package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/guillaumebreton/ruin/ofx"
	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/util"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import an ofx file into the ledger",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			util.Exit("Need to provide at least on directory/file as arg")
		}
		var total int = 0
		for _, arg := range args {
			dir, err := isDir(arg)
			if err != nil {
				util.HandleError(err, fmt.Sprintf("Directory %s does not exist", arg))
				continue
			}
			if dir {
				files, err := ioutil.ReadDir(arg)
				if err != nil {
					util.HandleError(err, fmt.Sprintf("Fail to read dir %s", arg))
					continue
				}
				for _, f := range files {
					if strings.HasSuffix(f.Name(), ".ofx") {
						path := path.Join(path.Dir(arg), f.Name())
						count, err := importFile(path)
						util.ExitOnError(err, fmt.Sprintf("Fail to read dir %s", arg))
						total += count
					}
				}
			} else {
				count, err := importFile(arg)
				util.ExitOnError(err, fmt.Sprintf("Fail to read dir %s", arg))
				total += count
			}
		}
		var tagged []*service.Transaction
		if ledger.Dirty {
			ledger.Reindex()
			tagged = ledger.Autotag()
			err := ledger.Save(ledgerFile)
			util.ExitOnError(err, "Fail to save ledger")
		}
		for _, v := range tagged {
			fmt.Printf("Transaction %d - %s (%.02f) was tagged %s\n", v.Number, v.Description, v.Amount, v.Category)
		}
		fmt.Printf("%d transaction(s) added\n", total)
	},
}

func importFile(file string) (int, error) {
	o, err := ofx.Parse(file)
	util.ExitOnError(err, "Fail to parse OFX file")
	count := 0
	for _, tx := range o.Transactions {
		a := tx.GetAmount()
		if ledger.Add(tx.ID, tx.GetDate(), tx.TxType, tx.Description, a) {
			count++
		}
	}
	if o.GetBalanceDate().After(ledger.BalanceDate) {
		ledger.Balance = o.Balance
		ledger.BalanceDate = o.GetBalanceDate()
		ledger.Dirty = true
	}

	if count > 0 {
		ledger.Dirty = true
	}

	return count, nil
}

func isDir(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return fi.Mode().IsDir(), nil
}

func init() {
	RootCmd.AddCommand(importCmd)
}
