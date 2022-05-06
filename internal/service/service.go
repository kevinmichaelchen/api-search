package service

import (
	"context"
	"github.com/kevinmichaelchen/api-search/internal/idl/coop/drivers/search/v1beta1"
	"github.com/meilisearch/meilisearch-go"
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

func driverToMap(in *v1beta1.Driver) map[string]interface{} {
	out := make(map[string]interface{})
	out["id"] = in.GetId()
	out["first_name"] = in.GetFirstName()
	out["last_name"] = in.GetLastName()
	out["email"] = in.GetEmail()
	out["phone"] = in.GetPhone()
	return out
}

func driverFromMap(in map[string]interface{}) *v1beta1.Driver {
	out := new(v1beta1.Driver)
	out.Id = in["id"].(string)
	out.FirstName = in["first_name"].(string)
	out.LastName = in["last_name"].(string)
	out.Email = in["email"].(string)
	out.Phone = in["phone"].(string)
	return out
}

func (s *Service) Ingest(ctx context.Context, req *v1beta1.IngestRequest) (*v1beta1.IngestResponse, error) {
	index := s.searchClient.Index(indexDrivers)
	var documents []map[string]interface{}
	for _, e := range req.GetPayloads() {
		if e.GetDriver() != nil {
			documents = append(documents, driverToMap(e.GetDriver()))
		}
	}
	_, err := index.AddDocuments(documents)
	if err != nil {
		return nil, err
	}
	return &v1beta1.IngestResponse{}, nil
}

func (s *Service) Search(ctx context.Context, req *v1beta1.SearchRequest) (*v1beta1.SearchResponse, error) {
	searchRes, err := s.searchClient.
		Index(indexDrivers).
		Search(req.GetQuery(),
			&meilisearch.SearchRequest{
				Limit: int64(req.GetLimit()),
			})
	if err != nil {
		return &v1beta1.SearchResponse{}, nil
	}
	var hits []*v1beta1.Payload
	for _, e := range searchRes.Hits {
		m, ok := e.(map[string]interface{})
		if ok {
			hits = append(hits, &v1beta1.Payload{
				Payload: &v1beta1.Payload_Driver{
					Driver: driverFromMap(m),
				},
			})
		} else {
			s.logger.Println("Type mismatch: Expected search hit to be map[string]interface{}")
		}
	}
	return &v1beta1.SearchResponse{
		Hits: hits,
	}, nil
}
