package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/guillaumebreton/ruin/util"
	"github.com/spf13/cobra"
)

// budgetsCmd represents the budgets command
var (
	budgetValue  string
	budgetDelete bool
)
var budgetsCmd = &cobra.Command{
	Use:   "budgets",
	Short: "Budgets management command",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {

			month := time.Now().Format("1")
			year := time.Now().Format("2006")
			period := fmt.Sprintf("%s-%s", month, year)
			budget := ledger.GetBudget(period)
			if budget == nil {
				fmt.Printf("Budget for month &d/%d doesn't exist, do you want to create it? [Y/N]\n")
				var response string
				_, err := fmt.Scanln(&response)
				if err != nil {
					log.Fatal(err)
				}
				if response == "Y" {
					budget = service.NewBudget()
					ledger.UpdateBudget(period, budget)
					ledger.Save(ledgerFile)
				} else {
					util.Exit("No budger found")
				}
			}
			RenderBudgetsListText(budget)
		}
		// else {
		// if !budgetDelete && budgetValue == "" {
		// 	util.Exit("Please use -v with a value or -d")
		// }
		// for _, budget := range args {
		// 	if budgetDelete {
		// 		ledger.DeleteBudget(budget)
		// 	} else {
		// 		var v float64
		// 		var err error
		// 		if budgetValue == "" {
		// 			v = 0
		// 		} else {
		// 			v, err = strconv.ParseFloat(budgetValue, 64)
		// 		}
		// 		if err != nil {
		// 			util.ExitOnError(err, "Fail to parese the value")
		// 		}
		// 		err = ledger.SetBudget(budget, v)
		// 		if err != nil {
		// 			util.ExitOnError(err, "Fail to set budget value")
		// 		}
		// 	}
		// }
		// ledger.Save(ledgerFile)
		// }
	},
}

func init() {
	RootCmd.AddCommand(budgetsCmd)
	budgetsCmd.Flags().StringVarP(&budgetValue, "value", "v", "", "Set the threshold value of the budget")
	budgetsCmd.Flags().BoolVarP(&budgetDelete, "delete", "d", false, "delete the budget")
}

func RenderBudgetsListText(budgets service.Budget) {
	table := table.NewTable()
	table.SetHeader("CATEGORY", "VALUE")
	var sum float64
	for k, v := range budgets {
		sum += v
		table.Append(k, v)
	}
	table.AppendSeparator()
	table.Append("TOTAL", sum)
	table.Render(os.Stdout)
}
