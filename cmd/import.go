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

	"github.com/anupcshan/ofx"
	"github.com/guillaumebreton/gobud/service"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		bs, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err: %v\n", err)
			os.Exit(1)
		}

		// Clean the files
		str := strings.Replace(string(bs), ",", ".", -1)
		re := regexp.MustCompile("(\\d)\\s(\\d)")
		str = re.ReplaceAllString(str, "${1}${2}")

		doc, err := ofx.Parse(strings.NewReader(str))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err: %v\n", err)
			os.Exit(1)
		}
		l, err := service.LoadLedger()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err: %v\n", err)
			os.Exit(1)
		}
		// TODO print the number of added task
		count := 0
		for _, tx := range doc.Transactions {
			a, _ := tx.Amount.Value.Float64()
			if l.Add(tx.ID, tx.PostedDate, tx.Type.String(), tx.Description, a) {
				count++
			}
		}
		fmt.Printf("%d transactions added\n", count)
		l.Save()
	},
}

func init() {
	txCmd.AddCommand(importCmd)
}
