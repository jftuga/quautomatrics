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
	"github.com/buger/jsonparser"
	"github.com/jftuga/quautomatrics/distribution"
	"github.com/jftuga/quautomatrics/library"
	mailingList "github.com/jftuga/quautomatrics/mailinglist"
	"github.com/jftuga/quautomatrics/survey"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
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

type DistributionEnvelope struct {
	libraryNameId string
	messageNameId string
	mailingListNameId string
	surveyNameId string
	sendDate string
	expirationDate string
	fromName string
	subject string
	replyToEmail string
	fromEmail string
}

var configFile, outputFile, sendDate, expirationDate string

func init() {
	rootCmd.AddCommand(createDistributionCmd)

	createDistributionCmd.Flags().StringVarP(&configFile, "configFile", "c", "", "JSON config file containing: fromName,fromEmail,replyToEmail,subject")
	createDistributionCmd.MarkFlagRequired("configFile")

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
	envelope := validateParameters()
	fmt.Println("envelope:", envelope)
	fmt.Println()
	dist := createDistributionEnvelope(envelope)
	fmt.Println(dist)
}

func createDistributionEnvelope(envelope *DistributionEnvelope) string {
	var dist string
	dist = distribution.JsonDistributionTemplate
	dist = strings.Replace(dist, "__LIBRARYID__", envelope.libraryNameId, 1)
	dist = strings.Replace(dist, "__MESSAGEID__", envelope.messageNameId, 1)
	dist = strings.Replace(dist, "__MAILINGLISTID__", envelope.mailingListNameId, 1)
	dist = strings.Replace(dist, "__FROMNAME__", envelope.fromName, 1)
	dist = strings.Replace(dist, "__REPLYTOEMAIL__", envelope.replyToEmail, 1)
	dist = strings.Replace(dist, "__FROMEMAIL__", envelope.fromEmail, 1)
	dist = strings.Replace(dist, "__SUBJECT__", envelope.subject, 1)
	dist = strings.Replace(dist, "__SURVEYID__", envelope.surveyNameId, 1)
	dist = strings.Replace(dist, "__EXPIRATIONDATE__", envelope.expirationDate, 1)
	dist = strings.Replace(dist, "__SENDDATE__", envelope.sendDate, 1)

	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error #16026: unable to open output file for writing: '%s'\n%s\n", outputFile, err)
	}
	f.WriteString(dist)
	f.Close()
	fmt.Printf("Saved distribution to file: '%s'\n", outputFile)

	return dist
}

func validateParameters() *DistributionEnvelope {
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

	// validate the opening of the output file
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error #40026: unable to open output file for writing: '%s'\n%s\n", outputFile, err)
	}
	f.Close()

	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error #26021: Unable to read config file: '%s'\n%s\n", configFile, err)
	}

	var fromName, replyToEmail, fromEmail, subject string
	handler := func(keyB []byte, valueB []byte, dataType jsonparser.ValueType, offset int) error {
		key := string(keyB)
		value:= string(valueB)

		if key == "fromName" {
			if len(value) < 5 {
				log.Fatalf("Error #26051: Config file has invalid %s: '%s'\n", key, value)
			}
			fromName = value
		} else if key == "subject" {
			if len(value) < 5 {
				log.Fatalf("Error #26052: Config file has invalid %s: '%s'\n", key, value)
			}
			subject = value
		}  else if key == "replyToEmail" {
			if isValidEmail(value) == false {
				log.Fatalf("Error #26053: Config file has invalid %s: '%s'\n", key, value)
			}
			replyToEmail = value
		} else if key == "fromEmail" {
			if isValidEmail(value) == false {
				log.Fatalf("Error #26054: Config file has invalid %s: '%s'\n", key, value)
			}
			fromEmail = value
		}
		return nil
	}
	jsonparser.ObjectEach(config, handler)
	//fmt.Println("1:", libraryNameId, messageNameId, mailingListNameId, surveyNameId, sendDate, expirationDate)
	//fmt.Println("2:", fromName, subject, replyToEmail, fromEmail)

	var envelope DistributionEnvelope
	envelope.libraryNameId = libraryNameId
	envelope.messageNameId = messageNameId
	envelope.mailingListNameId = mailingListNameId
	envelope.surveyNameId = surveyNameId
	envelope.sendDate = sendDate
	envelope.expirationDate = expirationDate
	envelope.fromName = fromName
	envelope.subject = subject
	envelope.replyToEmail = replyToEmail
	envelope.fromEmail = fromEmail

	return &envelope
}

func isValidEmail(address string) bool {
	// https://golangcode.com/validate-an-email-address/
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(address) < 3 && len(address) > 254 {
		return false
	}
	return emailRegex.MatchString(address)
}
