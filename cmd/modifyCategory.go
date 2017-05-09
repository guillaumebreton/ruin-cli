package cmd

import (
	"github.com/guillaumebreton/ruin/util"
	"github.com/spf13/cobra"
)

// renameCategoryCmd represents the renameCategory command
var modifyCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Modify a category",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			util.Exit("Please provide a category name and its new name")
		}
		ledger.RenameCategory(args[0], args[1])
		ledger.Save(ledgerFile)
	},
}

func init() {
	modifyCmd.AddCommand(modifyCategoryCmd)
}
