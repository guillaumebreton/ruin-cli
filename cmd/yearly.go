package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/guillaumebreton/ruin/util"
	"github.com/jinzhu/now"
	"github.com/spf13/cobra"
)

// yearlyCmd represents the yearly command
var yearlyCmd = &cobra.Command{
	Use:   "yearly",
	Short: "Generate a yearly report",
	Run: func(cmd *cobra.Command, args []string) {

		month, err := strconv.Atoi(time.Now().Format("1"))
		if err != nil {
			util.ExitOnError(err, "Fail to compute current month")
		}

		data := map[string][]float64{}

		soy := now.BeginningOfYear()

		f := service.NewFilter()
		f.StartDate = soy
		f.EndDate = time.Now()
		txs := ledger.GetTransactions(f)

		for _, tx := range txs {
			arr, ok := data[tx.Category]
			if !ok {
				arr = make([]float64, month+1)
			}
			txMonth, err := strconv.Atoi(tx.UserDate.Format("1"))
			if err != nil {
				util.ExitOnError(err, "Fail to parse tx date")
			}
			arr[txMonth-1] = arr[txMonth-1] + tx.Amount
			arr[len(arr)-1] = arr[len(arr)-1] + tx.Amount
			data[tx.Category] = arr
		}
		renderYearlyReport(month, data)
	},
}

func renderYearlyReport(month int, data map[string][]float64) {

	table := table.NewTable()
	header := make([]string, 2+month)
	header[0] = "CATEGORY"
	header[len(header)-1] = "TOTAL"
	for i := 1; i <= month; i++ {
		header[i] = strings.ToUpper(fmt.Sprintf("%s", time.Month(i)))
	}
	table.SetHeader(header)

	//the last row
	total := make([]float64, month+1)

	// sort the categories
	keys := make([]string, 0, len(data))
	for k, _ := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		vs := data[k]
		row := make([]string, 2+month)
		row[0] = k
		for i, v := range vs {
			row[i+1] = fmt.Sprintf("%0.2f", v)
			total[i] = total[i] + v
		}
		table.Append(row)
	}
	totalRow := make([]string, 2+month)
	totalRow[0] = "TOTAL"
	for k, v := range total {
		totalRow[k+1] = fmt.Sprintf("%0.2f", v)
	}
	table.AppendSeparator()
	table.Append(totalRow)
	table.Render(os.Stdout)

}

func init() {
	RootCmd.AddCommand(yearlyCmd)
}
