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
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/jftuga/quautomatrics/rest"
	"github.com/jftuga/quautomatrics/survey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"time"
)

var exportSurveyName string

// exportResponsesCmd represents the exportResponses command
var exportResponsesCmd = &cobra.Command{
	Use:   "exportResponses",
	Short: "Export survey responses in CSV format",
	Run: func(cmd *cobra.Command, args []string) {
		exportResponses()
	},
}

func init() {
	rootCmd.AddCommand(exportResponsesCmd)
	exportResponsesCmd.Flags().StringVarP(&exportSurveyName, "survey", "s", "", "preexisting survey name")
}

// extractCSVFileFromZip - The Qualtrics artifact is a zip data containing the survey name
//   example: survey="My Survey", then the packed file="My Survey.csv"
func extractCSVFileFromZip(zipData []byte, csvName string) string {
	/*
		// use this code to extract from a file instead of a stream...
		r, err := zip.OpenReader(zipFile)
		if err != nil {
			log.Fatalf("Error #34803: Unable to unzip data for survey='%s'\n%s\n", csvName, err)
		}
		defer r.Close()
	*/

	r, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		log.Fatalf("Error #39045: Unable to unzip data for survey='%s'\n%s\n", csvName, err)
	}

	csvName = csvName + ".csv"
	var packedFile *zip.File
	found := false
	for _, f := range r.File {
		if f.Name == csvName {
			packedFile = f
			found = true
			break
		}
	}

	if !found {
		log.Fatalf("Error #98432: Zip file did not contain survey='%s'\n", csvName)
	}

	outFile, err := os.OpenFile(csvName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Error #98902: Unable to open file for writing: '%s'\n'%s'\n", csvName, err)
	}

	rc, err := packedFile.Open()
	if rc == nil {
		log.Fatalf("Error #98903: Unable to packed file for reading: '%s'\n'%s'\n", csvName, err)
	}
	if err != nil {
		log.Fatalf("Error #98904: Unable to packed file for reading: '%s'\n'%s'\n", csvName, err)
	}

	_, err = io.Copy(outFile, rc)
	if err != nil {
		log.Fatalf("Error #98906: I/O Copy error with packed file: '%s'\n'%s'\n", csvName, err)
	}

	err = outFile.Close()
	if err != nil {
		log.Fatalf("Error #98907: Unable to close packed file: '%s'\n'%s'\n", csvName, err)
	}

	err = rc.Close()
	if err != nil {
		log.Fatalf("Error #98908: Unable to close packed file: '%s'\n'%s'\n", csvName, err)
	}

	return csvName
}

func statusCheck(surveyID, progressID string) []string {
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)

	allObjects := r.Get(fmt.Sprintf("/surveys/%s/export-responses/%s", surveyID, progressID) )
	if len(allObjects) == 0 {
		log.Fatalf("Error #488344: Status check failed for survey='%s', progressID='%s'\n", surveyID, progressID)
	}

	var status string
	result, _, _, err := jsonparser.Get([]byte(allObjects), "result")
	if err != nil {
		log.Fatalf("Error #22132: parsing JSON for key='result'\n%s\n", result)
	}

	status, err = jsonparser.GetString(result, "status")
	if err != nil {
		log.Fatalf("Error #22953: Unable to get status: '%s'\n", err)
	}

	var fileId string
	fileId, _ = jsonparser.GetString(result, "fileId")

	return []string { status, fileId}
}

// exportResponses - save the survey responses to the given CSV file
// see also: https://api.qualtrics.com/instructions/reference/responseImportsExports.json
func exportResponses() {
	if len(exportSurveyName) == 0 {
		log.Fatalln("Error #98004: The -s option expects a survey name")
	}

	exportSurvey := survey.New(exportSurveyName)

	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)

	payLoad := `{"format": "csv"}`
	allObjects := r.Post(fmt.Sprintf("/surveys/%s/export-responses", exportSurvey.Id), []byte(payLoad))
	if len(allObjects) == 0 {
		log.Fatalf("Error #56883: No exports were returned for exportSurvey='%s'\n", exportSurveyName)
	}

	// issue the export request
	result, _, _, err := jsonparser.Get([]byte(allObjects), "result")
	if err != nil {
		log.Fatalf("Error #56892: parsing JSON for key='result'\n%s\n", result)
	}

	progressId, err := jsonparser.GetString(result, "progressId")
	if err != nil {
		log.Fatalf("Error #34886: Unable to get progressId: '%s'\n", err)
	}

	// loop until the report generation is complete
	var status, fileId string
	count := 0
	for {
		blob := statusCheck(exportSurvey.Id, progressId)
		status = blob[0]
		fileId = blob[1]
		if status == "complete" {
			break
		}
		count += 1
		fmt.Printf("[%d] waiting for survey export completion...\n", count)
		time.Sleep( 1 * time.Second)
	}

	// download the report
	zipFile := r.Get(fmt.Sprintf("/surveys/%s/export-responses/%s/file", exportSurvey.Id, fileId))
	if len(zipFile) == 0 {
		log.Fatalf("Error #73902: Unable to download report for exportSurvey='%s'\n", exportSurveyName)
	}

	// extract CSV file from Zip stream and then save it
	savedCSVFile := extractCSVFileFromZip([]byte(zipFile), exportSurveyName)
	fmt.Println("Saved CSV file:", savedCSVFile)
}
