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
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	mailingList "github.com/jftuga/quautomatrics/mailinglist"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

// createContactsCmd represents the createContacts command
var createContactsCmd = &cobra.Command{
	Use:   "createContacts",
	Short: "Add contacts to a mailing list",
	Long:  `This will add contacts to an existing list. It does not check for duplicates.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("createContacts called")
		mList := mailingList.New(mailingListName)
		newContacts := getCSVEntries()
		CreateContacts(mList, newContacts)
	},
}

func init() {
	rootCmd.AddCommand(createContactsCmd)
	createContactsCmd.Flags().StringVarP(&mailingListName, "mailingList", "m", "", "name of qualtrics mailing list")
	createContactsCmd.MarkFlagRequired("mailingList")
	createContactsCmd.Flags().StringVarP(&csvFile, "csvFile", "c", "", "csv filename; format: first,last,email")
	createContactsCmd.MarkFlagRequired("csvFile")
}

// CreateContacts - create a group of new contacts for this mailing list
// new contacts are located in the CSV file given on command line
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
			log.Fatalf("Error #90774: parsing JSON for key='meta'\n%s\n", meta)
		}
	}
}
