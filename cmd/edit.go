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
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/guillaumebreton/ruin/service"
	"github.com/guillaumebreton/ruin/table"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a budget",
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
		name := args[0]

		l, err := service.LoadLedger()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err: %v\n", err)
			os.Exit(1)
		}
		// create a temp file
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			// handle error
			fmt.Fprintf(os.Stderr, "Err: %v\n", err)
			os.Exit(1)
		}

		defer os.Remove(tmpfile.Name())
		contentHash, content := GenerateData(l)

		if _, err := tmpfile.Write([]byte(content)); err != nil {
			// handle error
			fmt.Fprintf(os.Stderr, "Err: %v\n", err)
			os.Exit(1)
		}
		editor := os.Getenv("EDITOR")
		path, err := exec.LookPath(editor)

		if err != nil {
			fmt.Printf("Error %s while looking up for <%s>!!", path, editor)
		}

		editorCmd := exec.Command(path, tmpfile.Name())
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		err = editorCmd.Start()
		if err != nil {
			fmt.Printf("Start failed: %s", err)
		}
		err = editorCmd.Wait()
		if err != nil {
			fmt.Printf("Run failed: %s", err)
		} else {
			//load the file and check the
			s, err := ioutil.ReadFile(tmpfile.Name())
			if err != nil {
				fmt.Printf("Read failed: %s", err)
				os.Exit(1)
			}

			hasher := md5.New()
			sh := hex.EncodeToString(hasher.Sum([]byte(s)))
			if sh == contentHash {
				fmt.Println("Nothing changed")
				os.Exit(1)
			}

			file, err := os.Open(tmpfile.Name())
			if err != nil {
				fmt.Printf("Read failed: %s", err)
				os.Exit(1)
			}
			defer file.Close()

			budget := service.NewBudget()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				t := scanner.Text()
				if t == "" {
					continue
				}
				arr := strings.Split(t, ":")
				if len(arr) != 2 {
					fmt.Printf("Read failed: %s", err)
					os.Exit(1)
				}
				category := strings.TrimSpace(arr[0])
				v, err := strconv.ParseFloat(strings.TrimSpace(arr[1]), 64)
				if err != nil {
					fmt.Printf("Parse failed: %s", err)
					os.Exit(1)
				}
				budget.Set(category, v)
			}
			l.Budgets[name] = budget
			l.Save()
		}
	},
}

func init() {
	RootCmd.AddCommand(editCmd)

}

func GenerateData(l *service.Ledger) (string, string) {
	buf := bytes.NewBufferString("")
	table := table.NewTable()
	table.SetColumnSeparator(":     ")
	table.SetRowSeparator("")
	table.SetCenterSeparator("")
	table.SetBorder(false)

	for k, v := range l.Budgets {
		table.Append([]string{k, fmt.Sprintf("%.2f", v)})
	}
	table.Render(buf) // Send output
	content := buf.String()
	hasher := md5.New()
	h := hex.EncodeToString(hasher.Sum([]byte(content)))
	return h, content
}
