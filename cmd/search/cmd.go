package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "search",
	Short: "Search is a tool to make gRPC requests",
	Long: `Search is a tool to make gRPC requests built with
                love by The Drivers Coop.`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var conn *grpc.ClientConn

var query string
var limit int32
var ingestPath string
var genPath string

func init() {
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(ingestCmd)
	rootCmd.AddCommand(generateCmd)

	ingestCmd.AddCommand(ingestDriversCmd)

	searchCmd.PersistentFlags().Int32VarP(&limit, "limit", "l", 10, "Max number of results to return")
	searchCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "Search query")

	generateCmd.PersistentFlags().StringVarP(&genPath, "file", "f", "fake-drivers.csv", "Path of CSV file")
	ingestCmd.PersistentFlags().StringVarP(&ingestPath, "file", "f", "seed-drivers.json", "Path of JSON import file")

	log.Println("Initializing gRPC connection...")
	var err error
	conn, err = grpc.Dial("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial gRPC connection: %v", err)
	}
	log.Println("Initialized gRPC connection.")
}

func marshalProto(m proto.Message) (string, error) {
	b, err := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
