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
	"github.com/jftuga/quautomatrics/library"
	mailingList "github.com/jftuga/quautomatrics/mailinglist"
	"github.com/jftuga/quautomatrics/survey"
	"github.com/spf13/cobra"
	"log"
	"regexp"
)

// createDistributionCmd represents the createDistribution command
var createDistributionCmd = &cobra.Command{
	Use:   "createDistribution",
	Short: "Create a distribution file in JSON format",
	Long: `Long desc`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("createDistribution called")
		createDistribution()
	},
}

var outputFile, sendDate, expirationDate string

func init() {
	rootCmd.AddCommand(createDistributionCmd)
	createDistributionCmd.Flags().StringVarP(&outputFile, "outputFile", "o", "", "save JSON data to this file (will overwrite existing file)")
	createDistributionCmd.MarkFlagRequired("outputFile")

	createDistributionCmd.Flags().StringVarP(&libraryName, "libraryName", "l", "", "name of existing library")
	createDistributionCmd.MarkFlagRequired("libraryName")

	createDistributionCmd.Flags().StringVarP(&messageName, "messageName", "m", "", "name of existing library message")
	createDistributionCmd.MarkFlagRequired("messageName")

	createDistributionCmd.Flags().StringVarP(&mailingListName, "mailingListName", "n", "", "name of existing mailing list")
	createDistributionCmd.MarkFlagRequired("mailingListName")

	createDistributionCmd.Flags().StringVarP(&surveyName, "surveyName", "s", "", "name of existing survey")
	createDistributionCmd.MarkFlagRequired("surveyName")

	createDistributionCmd.Flags().StringVarP(&sendDate, "sendDate", "d", "", "format: YYYY-MM-DDTHH:MM:SSZ")
	createDistributionCmd.MarkFlagRequired("sendDate")

	createDistributionCmd.Flags().StringVarP(&expirationDate, "expirationDate", "e", "", "example: 2021-01-03T18:05:30Z")
	createDistributionCmd.MarkFlagRequired("expirationDate")
}

func createDistribution() {
	validateParameters()
}

func validateParameters() {

	// libraryID
	lib := library.New(libraryName)
	libraryNameId := lib.Id

	// libraryMessageId
	messageNameId := ""
	if len(messageName) > 0 {
		messageNameId = lib.GetLibraryMessage(messageName)
		if len(messageNameId) < 0 {
			log.Fatalf("Error #40010: unable to get library message name: '%s'\n", messageName)
		}
	}

	// mailingListId
	mList := mailingList.New(mailingListName)
	mailingListNameId := mList.Id

	// surveyId
	survey := survey.New(surveyName)
	surveyNameId := survey.Id

	// https://stackoverflow.com/a/28022901
	validDateRE := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:Z|[+-][01]\\d:[0-5]\\d)$"
	validDate := regexp.MustCompile(validDateRE)
	var valid bool
	//var err error
	valid = validDate.Match([]byte(sendDate))
	if valid == false {
		log.Fatalf("Error #40020: Invalid sendDate: '%s'\n", sendDate)
	}
	valid = validDate.Match([]byte(expirationDate))
	if valid == false {
		log.Fatalf("Error #40021: Invalid expirationDate: '%s'\n", expirationDate)
	}

	// TODO: validate the opening of the output file

	fmt.Println(libraryNameId, messageNameId, mailingListNameId, surveyNameId, sendDate, expirationDate)

}