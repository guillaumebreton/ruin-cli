package cmd

import (
	"github.com/spf13/cobra"
)

// reindexCmd represents the reindex command
var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		ledger.Save(ledgerFile)
	},
}

func init() {
	RootCmd.AddCommand(reindexCmd)
}
