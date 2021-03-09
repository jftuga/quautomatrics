package cmd

import (
	"encoding/csv"
	mailingList "github.com/jftuga/quautomatrics/mailinglist"
	"io"
	"log"
	"os"
)

// getCSVEntries - convert a CSV entry into an array of Contact
// CSV format: first,last,email
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
			log.Printf("Warning #72469: unable to process CSV entry: %s, line: %d\n%s\n", csvFile, line, err)
		}
		allCSVEntries = append(allCSVEntries, mailingList.Contact{"", record[2], record[0], record[1]})
	}
	return allCSVEntries
}
