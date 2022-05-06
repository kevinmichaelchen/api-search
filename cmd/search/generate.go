package main

import (
	"encoding/csv"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var size int

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate CSV with fake data",
	Long:  `Generate CSV with fake data`,
	Run:   generateFakeCSV,
}

func generateFakeCSV(cmd *cobra.Command, args []string) {
	f, err := os.Create(genPath)
	if err != nil {
		log.Fatalf("failed to create/open file: %v", err)
	}
	w := csv.NewWriter(f)
	err = w.Write([]string{
		"id", "first_name", "last_name", "email", "phone_number",
	})
	if err != nil {
		log.Fatalf("failed to write csv headers: %v", err)
	}
	for i := 0; i < size; i++ {
		firstName := faker.FirstName()
		lastName := faker.LastName()
		err := w.Write([]string{
			xid.New().String(),
			firstName,
			lastName,
			fmt.Sprintf("%s.%s@%s", firstName, lastName, faker.DomainName()),
			faker.Phonenumber(),
		})
		if err != nil {
			log.Fatalf("failed to write record: %v", err)
		}
	}
	w.Flush()
}
