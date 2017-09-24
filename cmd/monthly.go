package cmd

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/guillaumebreton/ruin/util"
	"github.com/jinzhu/now"
	"github.com/spf13/cobra"
)

var monthlyReportMonth string
var monthlyReportYear string
var reportWithTransactions bool

type ReportBudget struct {
	Category     string
	Value        float64
	Transactions service.Transactions
}

// monthlyCmd represents the monthly command
var monthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Generate a monthly report",
	Run: func(cmd *cobra.Command, args []string) {
		if monthlyReportMonth == "" {
			monthlyReportMonth = time.Now().Format("Jan")
		}
		if monthlyReportYear == "" {
			monthlyReportYear = time.Now().Format("2006")
		}
		t, err := time.Parse("Jan 2006", fmt.Sprintf("%s %s", monthlyReportMonth, monthlyReportYear))
		if err != nil {
			util.ExitOnError(err, "Invalid date format")
		}
		som := now.New(t).BeginningOfMonth()
		eom := now.New(t).EndOfMonth()
		f := service.NewFilter()
		f.StartDate = som
		f.EndDate = eom

		txs := ledger.GetTransactions(f)

		// Get budgets
		budgets := ledger.GetBudget(t.Format("1-2006"))

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
	t := table.NewTable()
	t.Renderer[4] = table.PositiveRed
	t.Renderer[3] = table.RedGreen
	t.SetHeader("CATEGORY", "CURRENT", "BUDGETED", "LEFT", "OVERSPENT", "FUTURE")

	// the keys
	keys := make([]string, 0, len(report))
	for k, _ := range report {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var eom, overspentEom, sumBudgeted, totalLeft float64
	for _, k := range keys {
		v := report[k]
		var sum float64
		var overspent float64
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

		if v.Value == 0 && current < 0 {
			overspent = -1 * sum
		} else if v.Value != 0 && left < 0 {
			overspent = -1 * left
		}
		sumBudgeted += v.Value
		totalLeft += left

		t.Append(v.Category, current, v.Value, left, overspent, future)
		if left > 0 {
			eom += left
		}
		overspentEom += overspent
	}
	t.AppendSeparator()
	t.Append("BALANCE", balance, sumBudgeted, totalLeft, overspentEom, balance-eom)
	t.Render(os.Stdout)
}

func format(v float64) string {
	return fmt.Sprintf("%0.2f", v)
}

func RenderReportWithTransactions(report map[string]ReportBudget) {
	table := table.NewTable()
	table.SetHeader("#", "CATEGORY", "SPENT")

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
			table.Append("", strings.ToUpper(v.Category), "")
			table.AppendSeparator()
			for _, t := range v.Transactions {
				table.Append(fmt.Sprintf("%d", t.Number), t.Description, fmt.Sprintf("%0.2f", t.Amount))
			}
			table.AppendSeparator()
			table.Append("", "TOTAL", fmt.Sprintf("%0.2f", sum))
		}
	}

	table.Render(os.Stdout)
}

func init() {
	RootCmd.AddCommand(monthlyCmd)
	monthlyCmd.Flags().StringVarP(&monthlyReportMonth, "month", "m", "", "the report month (default to the current month)")
	monthlyCmd.Flags().StringVarP(&monthlyReportYear, "year", "y", "", "the report year (default to the current year)")
	monthlyCmd.Flags().BoolVarP(&reportWithTransactions, "with-transactions", "t", false, "report budget with associated transactions")

}
