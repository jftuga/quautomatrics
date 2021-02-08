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
	"github.com/spf13/viper"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

// uploadDistributionCmd represents the uploadDistribution command
var uploadDistributionCmd = &cobra.Command{
	Use:   "uploadDistribution",
	Short: "Upload a distribution file",
	Long: `
Once uploaded, the distribution will be active and emails will 
be sent out starting at 'sendDate'`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("uploadDistribution called")
		uploadDistribution()
	},
}

var distFile string

func init() {
	rootCmd.AddCommand(uploadDistributionCmd)
	uploadDistributionCmd.Flags().StringVarP(&distFile, "distFile", "d", "", "JSON distribution file created with 'createDistribution'")
	uploadDistributionCmd.MarkFlagRequired("distFile")
}

// uploadDistribution - this is the main entry point for this file
// given a JSON file created by the createDistribution' cli option,
// upload it and make this distribution available for sending out survy invitations
func uploadDistribution() {
	dist, err := ioutil.ReadFile(distFile)
	if err != nil {
		log.Fatalf("Error #50092: unable to read distFile: '%s'\n%s\n", distFile, err)
	}
	token := viper.GetString("X-API-TOKEN")
	dc := viper.GetString("DATACENTER")
	r := rest.New(token, dc)

	status := r.Post("/distributions", dist)
	fmt.Println("uploadDistribution status:", status)
}