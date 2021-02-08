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
	"github.com/jftuga/quautomatrics/library"
	mailingList "github.com/jftuga/quautomatrics/mailinglist"
	"github.com/jftuga/quautomatrics/survey"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// createDistributionCmd represents the createDistribution command
var createDistributionCmd = &cobra.Command{
	Use:   "createDistribution",
	Short: "Create a distribution file in JSON format",
	Long:  `
Macros can be used for date-time:

_NOW_       => replaced with current date/time such as 2006-01-02T15:04:05Z
_TODAY_     => replaced with current date such as 2006-01-02
_YMD_       => same as _TODAY_
_HMS_       => replaced with current time such as 15:04:05
_TOMORROW_  => replaced with tomorrow's date such as 2006-01-03'
_DAYS:n_    => replaced with n days in the future; when n=3 then 2006-01-05
`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("createDistribution called")
		createDistribution()
	},
}

type DistributionEnvelope struct {
	libraryNameId     string
	messageNameId     string
	mailingListNameId string
	surveyNameId      string
	sendDate          string
	expirationDate    string
	fromName          string
	emailSubject           string
	replyToEmail      string
	fromEmail         string
}

var configFile, outputFile, sendDate, expirationDate, emailSubject string

func init() {
	rootCmd.AddCommand(createDistributionCmd)

	createDistributionCmd.Flags().StringVarP(&configFile, "configFile", "c", "", "JSON config file containing: fromName,fromEmail,replyToEmail,emailSubject")
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

	createDistributionCmd.Flags().StringVarP(&emailSubject, "emailSubject", "j", "", "email subject")
	createDistributionCmd.MarkFlagRequired("emailSubject")

	createDistributionCmd.Flags().StringVarP(&sendDate, "sendDate", "d", "", "format: YYYY-MM-DDTHH:MM:SSZ")
	createDistributionCmd.MarkFlagRequired("sendDate")

	createDistributionCmd.Flags().StringVarP(&expirationDate, "expirationDate", "e", "", "example: 2021-01-03T18:05:30Z")
	createDistributionCmd.MarkFlagRequired("expirationDate")
}

// createDistribution - this is the main entry point for this file
func createDistribution() {
	envelope := validateParameters()
	createDistributionEnvelope(envelope)
}

// createDistributionEnvelope - create a JSON string suitable for the "/distributions" API endpoint
// the resulting string will be saved to a file given by this cmd-line option: --outputFile
func createDistributionEnvelope(envelope *DistributionEnvelope) string {
	var dist string
	dist = JsonDistributionTemplate
	dist = strings.Replace(dist, "__LIBRARYID__", envelope.libraryNameId, 1)
	dist = strings.Replace(dist, "__MESSAGEID__", envelope.messageNameId, 1)
	dist = strings.Replace(dist, "__MAILINGLISTID__", envelope.mailingListNameId, 1)
	dist = strings.Replace(dist, "__FROMNAME__", envelope.fromName, 1)
	dist = strings.Replace(dist, "__REPLYTOEMAIL__", envelope.replyToEmail, 1)
	dist = strings.Replace(dist, "__FROMEMAIL__", envelope.fromEmail, 1)
	dist = strings.Replace(dist, "__EMAILSUBJECT__", envelope.emailSubject, 1)
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

// extractTimeDuration - extract an integer value from t when given re
// example input: _DAYS:5_
// example output: 5
func extractTimeDuration(t string, re *regexp.Regexp) int {
	slots := re.FindStringSubmatch(t)
	if len(slots) != 2 {
		log.Fatalf("Error #89173: Invalid macro for regular expression: '%s\n", t)
	}
	duration, err := strconv.Atoi(slots[1])
	if err != nil {
		log.Fatalf("Error #89027: Invalid time duration for: '%s\n%s\n", t, err)
	}
	return duration
}

// translateTimeMacos - convert macros surrounded by underscores to their
// corresponding date / time values
// these are used by these cmd-line options: --sendDate and --expirationDate
func translateTimeMacos(t string) string {
	currentTime := time.Now().UTC()
	tomorrowTime := currentTime.AddDate(0, 0, 1)

	// https://www.golangprograms.com/get-current-date-and-time-in-various-format-in-golang.html
	t = strings.Replace(t, "_NOW_", currentTime.Format("2006-01-02T15:04:05Z"), 1)
	t = strings.Replace(t, "_TODAY_", currentTime.Format("2006-01-02"), 1)
	t = strings.Replace(t, "_YMD_", currentTime.Format("2006-01-02"), 1)
	t = strings.Replace(t, "_HMS_", currentTime.Format("15:04:05"), 1)
	t = strings.Replace(t, "_TOMORROW_", tomorrowTime.Format("2006-01-02"), 1)

	if strings.Contains(t, "_DAYS:") == true {
		daysRE, err := regexp.Compile("_DAYS:([0-9]+)_")
		if err != nil {
			log.Fatalf("Error #89052: Invalid macro for _DAYS:_: '%s\n%s\n", t, err)
		}
		duration := extractTimeDuration(t, daysRE)
		futureTime := currentTime.AddDate(0, 0, duration)
		t = daysRE.ReplaceAllString(t, futureTime.Format("2006-01-02"))
	}
	/*
		// this is buggy and unreliable...
		if strings.Contains(t,"_HOURS:") == true {
			hoursRE, err := regexp.Compile("_HOURS:([0-9]+)_")
			if err != nil {
				log.Fatalf("Error #89054: Invalid macro for _HOURS:_: '%s\n%s\n", t, err)
			}
			duration := extractTimeDuration(t, hoursRE)
			futureTime := currentTime.Add(time.Hour * time.Duration(duration))
			t = hoursRE.ReplaceAllString(t,futureTime.Format("15"))
		}
	*/
	return t

}

// validateParameters - Certify that all command-line parameter are valid
// Make API calls to ensure a valid: libraryID, mailingListId, surveyId
// perform date-time macro substitutions
// ensure that timestamps are in the correct format
// ensure the config file contains: fromName, replyToEmail, fromEmail
// ensure the output file can be opened for writing
// once everything is validated, return a DistributionEnvelope containing all of this information
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

	// emailSubject
	if len(emailSubject) <= 5 {
		log.Fatalf("Error #40015: email subject is too short: '%s'\n", emailSubject)
	}
	if len(emailSubject) >= 45 {
		log.Fatalf("Error #40016: email subject is too long: '%s'\n", emailSubject)
	}

	// mailingListId
	mList := mailingList.New(mailingListName)
	mailingListNameId := mList.Id

	// surveyId
	survey := survey.New(surveyName)
	surveyNameId := survey.Id

	// make time substitutions
	sendDate = translateTimeMacos(sendDate)
	expirationDate = translateTimeMacos(expirationDate)

	// validate datetime
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

	var fromName, replyToEmail, fromEmail string
	handler := func(keyB []byte, valueB []byte, dataType jsonparser.ValueType, offset int) error {
		key := string(keyB)
		value := string(valueB)

		if key == "fromName" {
			if len(value) < 5 {
				log.Fatalf("Error #26051: Config file has invalid %s: '%s'\n", key, value)
			}
			fromName = value
		} else if key == "replyToEmail" {
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
	//fmt.Println("2:", fromName, emailSubject, replyToEmail, fromEmail)

	var envelope DistributionEnvelope
	envelope.libraryNameId = libraryNameId
	envelope.messageNameId = messageNameId
	envelope.mailingListNameId = mailingListNameId
	envelope.surveyNameId = surveyNameId
	envelope.sendDate = sendDate
	envelope.expirationDate = expirationDate
	envelope.fromName = fromName
	envelope.emailSubject = emailSubject
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
