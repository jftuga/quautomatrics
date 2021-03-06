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
	"github.com/jftuga/quautomatrics/rest"
	"github.com/jftuga/quautomatrics/survey"
	"github.com/spf13/cobra"
)

// listSurveysCmd represents the listSurveys command
var listSurveysCmd = &cobra.Command{
	Use:   "listSurveys",
	Short: `Get a survey ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("listSurveys called")
		listSurveys()
	},
}

var surveyName string

func init() {
	rootCmd.AddCommand(listSurveysCmd)
	listSurveysCmd.Flags().StringVarP(&surveyName, "name", "n", "", "preexisting survey name")
	//listSurveysCmd.MarkFlagRequired("name")
}

// listSurveys - return the survey id when given a survey name with the --name cli option
func listSurveys() {
	if len(surveyName) == 0 {
		allSurveys := rest.GenericMap("/surveys", map[string]string {"name": "name","id": "id"})
		if len(allSurveys) == 0 {
			fmt.Println("No surveys were returned by the API.")
		}
		for name, id := range allSurveys {
			fmt.Println("surveyId:", id, " name:", name)
		}
		return
	}
	survey := survey.New(surveyName)
	fmt.Println("surveyId:", survey.Id)
}
