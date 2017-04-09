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
	"github.com/guillaumebreton/gobud/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var reportEndDate string
var reportStartDate string

type ReportBudget struct {
	Category     string
	Value        float64
	Transactions service.Transactions
}

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "A brief description of your command",
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
		if reportEndDate != "" {
			t, err := time.Parse("2006-01-02", reportEndDate)
			if err != nil {
				fmt.Println("Invalid end-date format")
				os.Exit(1)
			}
			f.EndDate = t
		}
		if reportStartDate != "" {
			t, err := time.Parse("2006-01-02", reportStartDate)
			if err != nil {
				fmt.Println("Invalid start-date format")
				os.Exit(1)
			}
			f.StartDate = t
		}
		txs := l.GetTransactions(f)

		// Get budgets
		c, _ := service.LoadConfig()
		budgets := c.GetBudgets()

		report := map[string]ReportBudget{}
		for c, v := range budgets {
			rb := ReportBudget{
				Category:     c,
				Value:        v,
				Transactions: []service.Transaction{},
			}
			for _, tx := range txs {
				if tx.Category == c {
					rb.Transactions = append(rb.Transactions, tx)
				}
			}
			report[c] = rb
		}

		RenderReport(report)
	},
}

func RenderReport(report map[string]ReportBudget) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Spent", "Reserved", "Left"})
	for _, v := range report {
		var sum float64 = 0
		for _, v := range v.Transactions {
			sum += v.Amount
		}
		table.Append([]string{v.Category, fmt.Sprintf("%0.2f", sum), fmt.Sprintf("%0.2f", v.Value), fmt.Sprintf("%0.2f", v.Value-sum)})
	}
	table.SetAutoMergeCells(true)
	table.SetAutoWrapText(false)
	table.Render() // Send output

}

func init() {
	RootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "the start date")
	reportCmd.Flags().StringVarP(&endDate, "end-date", "e", "", "the end date")

}
