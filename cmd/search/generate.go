package main

import (
	"encoding/csv"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/kevinmichaelchen/api-search/internal/service/driver"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
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
		driver.FieldID, driver.FieldFirstName, driver.FieldLastName,
		driver.FieldEmail, driver.FieldPhone,
		driver.FieldTLCNumber,
		//driver.FieldVehicleClass,
		//driver.FieldVehicleMake,
		//driver.FieldVehicleModel,
		//driver.FieldVehicleYear,
		//driver.FieldLicensePlate,
	})
	if err != nil {
		log.Fatalf("failed to write csv headers: %v", err)
	}
	for i := 0; i < size; i++ {
		firstName := faker.FirstName()
		lastName := faker.LastName()
		tlcNumber, err := faker.RandomInt(1111111, 9999999)
		if err != nil {
			log.Fatalf("failed to write csv headers: %v", err)
		}
		err = w.Write([]string{
			xid.New().String(),
			firstName,
			lastName,
			fmt.Sprintf("%s.%s@%s", firstName, lastName, faker.DomainName()),
			faker.Phonenumber(),
			strconv.Itoa(tlcNumber[0]),
		})
		if err != nil {
			log.Fatalf("failed to write record: %v", err)
		}
	}
	w.Flush()
}
