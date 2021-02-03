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
	"fmt"
	"github.com/jftuga/quautomatrics/survey"

	"github.com/spf13/cobra"
)

// listMailingListsCmd represents the listMailingLists command
var listMailingListsCmd = &cobra.Command{
	Use:   "listMailingLists",
	Short: "Get a mailing-list ID",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("listMailingLists called")
		listMailingLists()
	},
}

var mListName string

func init() {
	rootCmd.AddCommand(listMailingListsCmd)
	listMailingListsCmd.Flags().StringVarP(&mListName, "name", "n", "", "preexisting mailing-list name")
	listMailingListsCmd.MarkFlagRequired("name")
}

func listMailingLists() {
	mList := survey.New(mListName)
	fmt.Println("surveyId:", mList.Id)
}
