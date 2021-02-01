/*
Copyright © 2021 John Taylor

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jftuga/quautomatrics/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mailingList, csvFile string

// replaceContactsCmd represents the replaceContacts command
var replaceContactsCmd = &cobra.Command{
	Use:   "replaceContacts",
	Short: "Replace all mailing list entries with CSV entries",
	Long: `(1) Delete all entries from the given mailing list
(2) Import all entries from the CSV file

CSV file format (no header line):
first name, last name, email address`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("replaceContacts called")
		replaceContacts(args)
	},
}

func init() {
	rootCmd.AddCommand(replaceContactsCmd)
	replaceContactsCmd.Flags().StringVarP(&mailingList, "mailingList","m", "", "name of qualtrics mailing list")
	replaceContactsCmd.MarkFlagRequired("mailingList")
	replaceContactsCmd.Flags().StringVarP(&csvFile, "csvFile","c", "", "csv filename; format: first,last,email")
	replaceContactsCmd.MarkFlagRequired("csvFile")
}

func replaceContacts(args []string) {
	fmt.Println(mailingList)
	fmt.Println(csvFile)
	fmt.Println(viper.Get("X-API-TOKEN"))
	conn := rest.New(viper.GetString("X-API-TOKEN"),viper.GetString("DATACENTER"))
	allMailingLists := conn.Get("mailinglists?offset=0")
	fmt.Println(allMailingLists)
	fmt.Println("is valid JSON:", json.Valid([]byte(allMailingLists)))
}
