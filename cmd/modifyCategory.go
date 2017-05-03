package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
)

// renameCategoryCmd represents the renameCategory command
var modifyCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Modify a category",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintf(os.Stderr, "Please provide a category name and its new name")
			os.Exit(1)
		}
		ledger.RenameCategory(args[0], args[1])
		ledger.Save(ledgerFile)
	},
}

func init() {
	modifyCmd.AddCommand(modifyCategoryCmd)
}
