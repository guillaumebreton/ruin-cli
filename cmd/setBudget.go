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
	"github.com/spf13/cobra"
	"strconv"
)

// setBudgetCmd represents the set command
var setBudgetCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if len(args) < 2 {
			fmt.Println("Need to provide a category and a value")
			return
		}
		key := args[0]
		value := args[1]
		c, _ := service.LoadConfig()
		v, _ := strconv.ParseFloat(value, 64)
		c.SetBudget(key, v)
		c.Save()
	},
}

func init() {
	budgetCmd.AddCommand(setBudgetCmd)
}