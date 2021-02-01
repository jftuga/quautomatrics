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
	"github.com/buger/jsonparser"
	mailingList "github.com/jftuga/quautomatrics/mailinglist"
	"github.com/jftuga/quautomatrics/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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
		replaceContacts(args)
	},
}

func init() {
	rootCmd.AddCommand(replaceContactsCmd)
	replaceContactsCmd.Flags().StringVarP(&mailingListName, "mailingList", "m", "", "name of qualtrics mailing list")
	replaceContactsCmd.MarkFlagRequired("mailingList")
	replaceContactsCmd.Flags().StringVarP(&csvFile, "csvFile", "c", "", "csv filename; format: first,last,email")
	replaceContactsCmd.MarkFlagRequired("csvFile")
}

func replaceContacts(args []string) {
	//fmt.Println(mailingListName)
	//fmt.Println(csvFile)
	//fmt.Println(viper.Get("X-API-TOKEN"))
	conn := rest.New(viper.GetString("X-API-TOKEN"), viper.GetString("DATACENTER"))
	allMailingLists := conn.Get("mailinglists?offset=0")
	fmt.Println(allMailingLists)
	if json.Valid([]byte(allMailingLists)) != true {
		log.Fatalf("Invalid JSON returned from API:\n%s\n", allMailingLists)
		return
	}
	result, _, _, err := jsonparser.Get([]byte(allMailingLists), "result")
	if err != nil {
		log.Fatalf("Error #46237: parsing JSON for key='result'\n%s\n", result)
	}
	temp := string(result)
	fmt.Println("temp:", temp)
	hasList := false
	var name, id string
	_, err = jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, err = jsonparser.GetString(value, "name")
		if err != nil {
			log.Fatalf("Error #46885: parsing JSON for key='name'\n%s\n", value)
		}
		if name == mailingListName {
			hasList = true
			id, err = jsonparser.GetString(value, "id")
			if err != nil {
				log.Fatalf("Error #46005: parsing JSON for key='id'\n%s\n", value)
			}
			return
		}
	}, "elements")
	if err != nil {
		log.Fatalf("Error #38932: %s", err)
	}

	if hasList == false {
		log.Fatalf("Error #46376: mailing list does not exist: %s\n", mailingListName)
	}

	//fmt.Println("id:", id)
	mList := mailingList.New(mailingListName, id)
	allContacts := mList.GetAllContacts()
	for _, contact := range allContacts {
		fmt.Println(contact.Email)
		success := mList.DeleteContact(contact.Id)
		fmt.Println("delete success:", success)
		fmt.Println()
	}

	//part 2 - import from csv

}
