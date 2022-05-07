package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/kevinmichaelchen/api-search/internal/idl/coop/drivers/search/v1beta1"
	"github.com/kevinmichaelchen/api-search/internal/service/driver"
	"github.com/meilisearch/meilisearch-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

const (
	indexDrivers = "drivers"
)

type Service struct {
	logger       *log.Logger
	searchClient *meilisearch.Client
}

func NewService(logger *log.Logger, searchClient *meilisearch.Client) *Service {
	return &Service{logger: logger, searchClient: searchClient}
}

func encode(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(id))
}

func decode(id string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func driverToMap(in *v1beta1.Driver) map[string]interface{} {
	out := make(map[string]interface{})
	out[driver.FieldID] = encode(in.GetId())
	out[driver.FieldFirstName] = in.GetFirstName()
	out[driver.FieldLastName] = in.GetLastName()
	out[driver.FieldEmail] = in.GetEmail()
	out[driver.FieldPhone] = in.GetPhone()
	return out
}

func driverFromMap(in map[string]interface{}) (*v1beta1.Driver, error) {
	out := new(v1beta1.Driver)
	id, err := decode(in[driver.FieldID].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to decode ID: %w", err)
	}
	out.Id = id
	out.FirstName = in[driver.FieldFirstName].(string)
	out.LastName = in[driver.FieldLastName].(string)
	out.Email = in[driver.FieldEmail].(string)
	out.Phone = in[driver.FieldPhone].(string)
	return out, nil
}

func (s *Service) Ingest(ctx context.Context, req *v1beta1.IngestRequest) (*v1beta1.IngestResponse, error) {
	index := s.searchClient.Index(indexDrivers)
	var documents []map[string]interface{}
	for _, e := range req.GetDrivers().GetDrivers() {
		if e != nil {
			documents = append(documents, driverToMap(e))
		}
	}
	s.logger.Printf("Ingesting %d documents\n", len(documents))
	if len(documents) > 0 {
		s.logger.Println("Ingesting 1st doc:", documents[0])
	}
	res, err := index.AddDocuments(documents, "id")
	if err != nil {
		return nil, err
	}
	s.logger.Println("response", res)
	res, err = s.searchClient.WaitForTask(res)
	if err != nil {
		return nil, err
	}
	s.logger.Println("response", res)
	return &v1beta1.IngestResponse{
		Uid:      res.UID,
		IndexUid: res.IndexUID,
		Status:   string(res.Status),
		TaskType: res.Type,
		//Duration:   res.Duration,
		EnqueuedAt: timestamppb.New(res.EnqueuedAt),
		StartedAt:  timestamppb.New(res.StartedAt),
		FinishedAt: timestamppb.New(res.FinishedAt),
	}, nil
}

func (s *Service) Query(ctx context.Context, req *v1beta1.QueryRequest) (*v1beta1.QueryResponse, error) {
	searchRes, err := s.searchClient.
		Index(indexDrivers).
		Search(req.GetQuery(),
			&meilisearch.SearchRequest{
				Limit: int64(req.GetLimit()),
			})
	if err != nil {
		return nil, err
	}
	var drivers []*v1beta1.Driver
	for _, e := range searchRes.Hits {
		m, ok := e.(map[string]interface{})
		if ok {
			d, err := driverFromMap(m)
			if err != nil {
				return nil, err
			}
			drivers = append(drivers, d)
		} else {
			s.logger.Println("Type mismatch: Expected search hit to be map[string]interface{}")
		}
	}
	return &v1beta1.QueryResponse{
		Response: &v1beta1.QueryResponse_Drivers{
			Drivers: &v1beta1.DriverResponse{
				Results: drivers,
			},
		},
	}, nil
}
