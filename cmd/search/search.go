package main

import (
	"context"
	"github.com/kevinmichaelchen/api-search/internal/idl/coop/drivers/search/v1beta1"
	"github.com/spf13/cobra"
	"log"
)

var searchCmd = &cobra.Command{
	Use:   "query",
	Short: "Get search results",
	Long:  `Get search results`,
	Run:   search,
}

func search(cmd *cobra.Command, args []string) {
	// Create request
	req := &v1beta1.SearchRequest{
		Query: query,
		Limit: limit,
	}

	// Execute request
	client := v1beta1.NewSearchServiceClient(conn)
	res, err := client.Search(context.Background(), req)
	if err != nil {
		log.Fatalf("gRPC request failed: %v", err)
	}

	// Print response
	s, err := marshalProto(res)
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}
	log.Println(s)
}
