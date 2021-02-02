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
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	mailingList "github.com/jftuga/quautomatrics/mailinglist"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
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

func ReplaceContacts() {
	mList := mailingList.New(mailingListName)
	DeleteAllContacts(mList)
	newContacts := getCSVEntries()
	CreateContacts(mList, newContacts)
}

func DeleteAllContacts(mList *mailingList.MailingList) {
	allContacts := mList.GetAllContacts()
	for _, contact := range allContacts {
		fmt.Println("removing contact:", contact.Email)
		success := mList.DeleteContact(contact.Id)
		fmt.Println("delete success:", success)
		fmt.Println()
	}
}

func CreateContacts(mList *mailingList.MailingList, newContacts []mailingList.Contact) {
	// we don't need the the 'Id' field from contact
	type contactSmall struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	path := fmt.Sprintf("/mailinglists/%s/contacts", mList.Id)
	for _, contact := range newContacts {
		jsonData := new(bytes.Buffer)
		c := contactSmall{contact.FirstName, contact.LastName, contact.Email}
		fmt.Println("creating contact:", c.Email)
		err := json.NewEncoder(jsonData).Encode(c)
		if err != nil {
			log.Fatalf("Error #47922: unable to encode to JSON: %s\n%s\n", c, err)
		}
		buf, _ := ioutil.ReadAll(jsonData)
		request := mList.Conn.Rest.Post(path, buf)
		meta, _, _, err := jsonparser.Get([]byte(request), "meta")
		if err != nil {
			log.Fatalf("Error #73639: parsing JSON for key='meta'\n%s\n", meta)
		}
	}
}

func getCSVEntries() []mailingList.Contact {
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Error #70055: unable to open ")
	}
	r := csv.NewReader(file)

	var allCSVEntries []mailingList.Contact
	line := 0
	for {
		line += 1
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error #70062: unable to process CSV entry: %s, line: %d\n%s\n", csvFile, line, err)
		}
		allCSVEntries = append(allCSVEntries, mailingList.Contact{"", record[2], record[0], record[1]})
	}
	return allCSVEntries
}
