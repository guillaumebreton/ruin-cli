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
	"github.com/guillaumebreton/ruin/table"
	"github.com/spf13/cobra"
	"os"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Need a name")
			os.Exit(1)
		}
		l, err := service.LoadLedger()
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %v\n", err)
			os.Exit(1)
		}
		budget, err := l.GetBudget(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Budget name doesn't exist")
			os.Exit(1)
		}
		RenderBudget(budget.Values)
	},
}

func RenderBudget(values map[string]float64) {
	table := table.NewTable()
	table.SetHeader([]string{"CATEGORY", "VALUE"})
	var sum float64 = 0
	for k, v := range values {
		sum += v
		table.Append([]string{k, fmt.Sprintf("%0.2f", v)})
	}
	table.AppendSeparator()
	table.Append([]string{"TOTAL", fmt.Sprintf("%0.2f", sum)})
	table.Render(os.Stdout)

}
func init() {
	RootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
