package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
)

// renameCategoryCmd represents the renameCategory command
var renameCategoryCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a category",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintf(os.Stderr, "Please provide a category name and its new name")
			os.Exit(1)
		}
		ledger.RenameCategory(args[0], args[1])
		ledger.Save()
	},
}

func init() {
	RootCmd.AddCommand(renameCategoryCmd)
}
