package cmd

import (
	"github.com/spf13/cobra"
)

// modifyCmd represents the modify command
var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify a budget, a transaction or a category",
}

func init() {
	RootCmd.AddCommand(modifyCmd)
}
