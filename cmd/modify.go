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
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	modifyCategory string
	modifyDate     string
)

// modifyCmd represents the modify command
var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify a transaction",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var d time.Time
		var err error
		if len(args) != 1 {
			fmt.Println("Please proive a list of transaction id")
			os.Exit(1)
		}
		if modifyDate != "" {
			d, err = time.Parse(service.ShortFormat, modifyDate)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid date : %s", modifyDate)
				os.Exit(1)
			}

		}
		l, err := service.LoadLedger()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err: %v\n", err)
			os.Exit(1)
		}
		arr := strings.Split(args[0], ",")
		for _, v := range arr {
			id, err := strconv.Atoi(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s not found", v)
			} else {
				tx, err := l.GetTransaction(id)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%d not found", id)
				} else {
					if modifyCategory != "" {
						tx.Category = modifyCategory
					}
					if modifyDate != "" {
						println("Set date")
						tx.UserDate = d
					}
					l.UpdateTransaction(id, tx)
				}
			}
		}
		fmt.Printf("Transaction %s Updated\n", args[0])
		l.Save()
	},
}

func init() {
	RootCmd.AddCommand(modifyCmd)
	modifyCmd.Flags().StringVarP(&modifyCategory, "category", "c", "", "Set the transaction category")
	modifyCmd.Flags().StringVarP(&modifyDate, "date", "d", "", "Set the transaction date (eg. 2017-02-03)")
}
