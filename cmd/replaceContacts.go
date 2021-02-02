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
	mailingList "github.com/jftuga/quautomatrics/mailinglist"
	"github.com/spf13/cobra"
)

var mailingListName, csvFile string

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
		ReplaceContacts()
	},
}

func init() {
	rootCmd.AddCommand(replaceContactsCmd)
	replaceContactsCmd.Flags().StringVarP(&mailingListName, "mailingList", "m", "", "name of qualtrics mailing list")
	replaceContactsCmd.MarkFlagRequired("mailingList")
	replaceContactsCmd.Flags().StringVarP(&csvFile, "csvFile", "c", "", "csv filename; format: first,last,email")
	replaceContactsCmd.MarkFlagRequired("csvFile")
}

// ReplaceContacts - main entry point of this file; delete all current contacts, then create new contacts
// this requires --mailingList and --csvFile
func ReplaceContacts() {
	mList := mailingList.New(mailingListName)
	DeleteAllContacts(mList)
	newContacts := getCSVEntries()
	CreateContacts(mList, newContacts)
}

