package cmd

import (
	"fmt"
	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/spf13/cobra"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

var reportEndDate string
var reportStartDate string
var reportWithTransactions bool

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
		txs := ledger.GetTransactions(f)

		// Get budgets
		budgets := ledger.GetBudgets()

		report := map[string]ReportBudget{}

		for _, t := range txs {
			rb, ok := report[t.Category]
			if ok {
				rb.Transactions = append(rb.Transactions, t)
				report[t.Category] = rb
			} else {
				rb = ReportBudget{
					Category:     t.Category,
					Value:        0,
					Transactions: []*service.Transaction{t},
				}
				report[t.Category] = rb
			}
		}

		for c, v := range budgets {
			rb, ok := report[c]
			if ok {
				rb.Value = v
				report[c] = rb
			} else {
				rb := ReportBudget{
					Category:     c,
					Value:        v,
					Transactions: []*service.Transaction{},
				}
				report[c] = rb
			}

		}

		if reportWithTransactions {
			RenderReportWithTransactions(report)
		} else {
			if len(txs) > 0 {
				RenderReport(txs[0].Balance, report)
			} else {

				RenderReport(0, report)
			}
		}
	},
}

func RenderReport(balance float64, report map[string]ReportBudget) {
	table := table.NewTable()
	table.SetHeader([]string{"CATEGORY", "CURRENT", "FUTURE", "STATUS"})

	// the keys
	keys := make([]string, 0, len(report))
	for k, _ := range report {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var eom float64
	for _, k := range keys {
		v := report[k]
		var sum float64 = 0
		for _, v := range v.Transactions {
			sum += v.Amount
		}
		var current, future, left float64
		if v.Value != 0 && (-1*v.Value) < sum {
			current = sum
			future = -1 * v.Value
			left = math.Abs(v.Value) - math.Abs(current)
		} else if v.Value != 0 && (-1*v.Value) > sum {
			current = sum
			future = sum
			left = math.Abs(v.Value) - math.Abs(current)
		} else {
			current = sum
			future = sum
			left = 0
		}
		table.Append([]string{v.Category, format(current), format(future), format(left)})
		if left > 0 {
			eom += left
		}
	}
	table.AppendSeparator()
	table.Append([]string{"BALANCE", format(balance), format(balance - eom), ""})
	table.Render(os.Stdout)
}

func format(v float64) string {
	return fmt.Sprintf("%0.2f", v)
}

func RenderReportWithTransactions(report map[string]ReportBudget) {
	table := table.NewTable()
	table.SetHeader([]string{"#", "CATEGORY", "SPENT"})

	// the keys
	keys := make([]string, 0, len(report))
	for k, _ := range report {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var spentTotal, reservedTotal, leftTotal float64
	for _, k := range keys {
		v := report[k]
		var sum float64 = 0
		for _, t := range v.Transactions {
			sum += t.Amount
		}
		if sum != 0 {

			spentTotal += sum
			reservedTotal += -1 * v.Value
			leftTotal += v.Value - math.Abs(sum)
			table.AppendSeparator()
			table.Append([]string{"", strings.ToUpper(v.Category), ""})
			table.AppendSeparator()
			for _, t := range v.Transactions {
				table.Append([]string{fmt.Sprintf("%d", t.Number), t.Description, fmt.Sprintf("%0.2f", t.Amount)})
			}
			table.AppendSeparator()
			table.Append([]string{"", "TOTAL", fmt.Sprintf("%0.2f", sum)})
		}
	}

	table.Render(os.Stdout)
}

func init() {
	RootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringVarP(&reportStartDate, "start-date", "s", "", "the start date")
	reportCmd.Flags().StringVarP(&reportEndDate, "end-date", "e", "", "the end date")
	reportCmd.Flags().BoolVarP(&reportWithTransactions, "with-transactions", "t", false, "report budget with associated transactions")

}
