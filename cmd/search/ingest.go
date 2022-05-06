package main

import (
	"context"
	"encoding/csv"
	"github.com/kevinmichaelchen/api-search/internal/idl/coop/drivers/search/v1beta1"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
)

var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest search documents",
	Long:  `Ingest search documents`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

var ingestDriversCmd = &cobra.Command{
	Use:   "drivers",
	Short: "Ingest drivers",
	Long:  `Ingest drivers`,
	Run:   ingestDrivers,
}

func ingestDrivers(cmd *cobra.Command, args []string) {
	// Open file
	f, err := os.Open(ingestPath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	var payloads []*v1beta1.Payload

	// Read CSV records
	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to read csv headers: %v", err)
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Failed to read csv record: %v", err)
		}
		payloads = append(payloads, &v1beta1.Payload{
			Payload: &v1beta1.Payload_Driver{
				Driver: &v1beta1.Driver{
					Id:        record[0],
					FirstName: record[1],
					LastName:  record[2],
					Email:     record[3],
					Phone:     record[4],
				},
			},
		})
	}

	// Create request
	req := &v1beta1.IngestRequest{Payloads: payloads}
	s, err := marshalProto(req)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}
	log.Println(s)

	// Execute request
	client := v1beta1.NewSearchServiceClient(conn)
	res, err := client.Ingest(context.Background(), req)
	if err != nil {
		log.Fatalf("gRPC request failed: %v", err)
	}

	// Print response
	s, err = marshalProto(res)
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}
	log.Println(s)
}
