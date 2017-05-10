package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/util"
	"github.com/spf13/cobra"
)

var ledgerFile string
var ledger *service.Ledger

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ruin",
	Short: "Generates budget data from ofx files",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initLedger)
	RootCmd.PersistentFlags().StringVarP(&ledgerFile, "file", "f", "$HOME/.ruin.json", "Ledger file")
}

func initLedger() {
	var err error
	ledger, err = service.LoadLedger(ledgerFile)
	if err != nil {
		if ledgerFile == "$HOME/.ruin.json" {
			ledger = service.NewLedger()
			usr, err := user.Current()
			util.ExitOnError(err, "Fail to obtain user home directory")
			ledgerFile = usr.HomeDir + "/.ruin.json"
			if _, err := os.Stat(ledgerFile); os.IsNotExist(err) {
				err = ledger.Save(ledgerFile)
				util.ExitOnError(err, "Fail to create initial file")
			}
			ledger, _ = service.LoadLedger(ledgerFile)

		} else {
			util.Exit("Fail to load the ledger")
		}
	}
}
