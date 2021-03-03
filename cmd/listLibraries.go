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
	"github.com/jftuga/quautomatrics/rest"
	"github.com/spf13/cobra"
	"log"
)

var libraryName, messageName string
var allMessages bool

// listLibrariesCmd represents the listLibraries command
var listLibrariesCmd = &cobra.Command{
	Use:   "listLibraries",
	Short: `List all libraries. A library is needed in order to create a Distribution.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("listLibraries called")
		listLibraries()
	},
}

func init() {
	rootCmd.AddCommand(listLibrariesCmd)
	listLibrariesCmd.Flags().StringVarP(&libraryName, "name", "n", "", "preexisting library name")
	//listLibrariesCmd.MarkFlagRequired("name")
	listLibrariesCmd.Flags().StringVarP(&messageName, "message", "m", "", "preexisting library message")
	listLibrariesCmd.Flags().BoolVarP(&allMessages, "allMessages", "M", false, "list all library messages")
}

// listLibraries - output the Id and MessageID returned via the API
// when given a library name and message  (--name, --message cli options)
func listLibraries() {
	if len(libraryName) == 0 {
		if allMessages {
			log.Fatalln("Error #70062: The -M option must be used with -n option.")
		}
		allLibraries := rest.GenericMap("/libraries", map[string]string {"name": "libraryName","id": "libraryId"})
		if len(allLibraries) == 0 {
			fmt.Println("No libraries were returned by the API.")
		}
		for name, id := range allLibraries {
			fmt.Println("libraryId:", id, " name:", name)
		}
		return
	}

	lib := library.New(libraryName)
	fmt.Println("libraryId:", lib.Id)

	if len(messageName) > 0 {
		messageId := lib.GetLibraryMessage(messageName)
		if len(messageId) > 0 {
			fmt.Println("libraryMessageId:", messageId)
		}
	} else if allMessages {
		allLibraryMessages :=  lib.GetAllLibraryMessage()
		if len(allLibraryMessages) > 0 {
			for msgName, msgID := range allLibraryMessages {
				fmt.Println("messageId:", msgID, " name:", msgName)
			}
		}
	}
}
