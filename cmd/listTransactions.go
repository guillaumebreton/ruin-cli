// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/guillaumebreton/ruin/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var endDate string
var startDate string

// listTransactionCmd represents the listTransaction command
var listTransactionsCmd = &cobra.Command{
	Use:   "tx",
	Short: "List transactions",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		l, err := service.LoadLedger()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err: %v\n", err)
			os.Exit(1)
		}
		f := service.NewFilter()
		if endDate != "" {
			t, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				fmt.Println("Invalid end-date format")
				os.Exit(1)
			}
			f.EndDate = t
		}
		if startDate != "" {
			t, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				fmt.Println("Invalid start-date format")
				os.Exit(1)
			}
			f.StartDate = t
		}

		RenderTransactionListText(l.GetTransactions(f))
	},
}

func RenderTransactionListText(transactions []service.Transaction) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Date", "Description", "Category", "Amount"})
	for _, v := range transactions {
		table.Append([]string{fmt.Sprintf("%d", v.Number), v.Date.Format("2006-01-02"), v.Description, v.Category, fmt.Sprintf("%.2f", v.Amount)})
	}
	// table.SetAutoMergeCells(true)
	table.SetAutoWrapText(false)
	table.Render() // Send output
}
func init() {
	RootCmd.AddCommand(listTransactionsCmd)

	listTransactionsCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "the start date")
	listTransactionsCmd.Flags().StringVarP(&endDate, "end-date", "e", "", "the end date")

}
