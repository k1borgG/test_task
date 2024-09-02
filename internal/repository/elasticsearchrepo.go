package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/k1borgG/test_task/internal/dto"
)

type ElasticsearchRepository interface {
	IndexProduct(ctx context.Context, product dto.ProductDTO) (string, error)
	SearchProducts(ctx context.Context, searchQuery map[string]interface{}) ([]dto.ProductDTO, error)
}

type EsProduct struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Brand       string      `json:"brand"`
	Model       string      `json:"model"`
	Price       string      `json:"price"`
	Coordinates Coordinates `json:"coordinates"`
	CreatedAt   struct {
		Seconds int64 `json:"seconds"`
		Nanos   int32 `json:"nanos"`
	} `json:"created_at"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type EsHit struct {
	Source EsProduct `json:"_source"`
}

type EsResponse struct {
	Hits struct {
		Hits []EsHit `json:"hits"`
	} `json:"hits"`
}

type elasticsearchRepository struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearchRepository(client *elasticsearch.Client, index string) ElasticsearchRepository {
	return &elasticsearchRepository{
		client: client,
		index:  index,
	}
}

func (r *elasticsearchRepository) IndexProduct(ctx context.Context, product dto.ProductDTO) (string, error) {
	productID := uuid.New().String()

	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Printf("Error marshalling product: %v", err)
		return "", err
	}

	req := esapi.IndexRequest{
		Index:      r.index,
		DocumentID: productID,
		Body:       strings.NewReader(string(productJSON)),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Printf("Error indexing product in Elasticsearch: %v", err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Closing body error: %v", err)
		}
	}(res.Body)

	if res.IsError() {
		log.Printf("Error indexing product in Elasticsearch: %s", res.String())
		return "", fmt.Errorf("error indexing product: %s", res.String())
	}

	return productID, nil
}

func (r *elasticsearchRepository) SearchProducts(ctx context.Context, searchQuery map[string]interface{}) ([]dto.ProductDTO, error) {
	searchQueryJSON, err := json.Marshal(searchQuery)
	if err != nil {
		log.Printf("Error marshalling search query: %s", err)
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{"my_custom_index5"},
		Body:  strings.NewReader(string(searchQueryJSON)),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Fatalf("Error searching in Elasticsearch: %s", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Closing body error: %v", err)
		}
	}(res.Body)

	if res.IsError() {
		log.Printf("Error searching in Elasticsearch: %s", res.String())
		return nil, err
	}

	var esResp EsResponse
	if err = json.NewDecoder(res.Body).Decode(&esResp); err != nil {
		log.Printf("Error decode response from ES: %v", err)
		return nil, err
	}

	var products []dto.ProductDTO

	for _, hit := range esResp.Hits.Hits {
		product := dto.ProductDTO{
			Name:        hit.Source.Name,
			Description: hit.Source.Description,
			Brand:       hit.Source.Brand,
			Model:       hit.Source.Model,
			Price:       hit.Source.Price,
			Coordinates: dto.CoordinatesDTO{
				Lat: hit.Source.Coordinates.Lat,
				Lon: hit.Source.Coordinates.Lon,
			},
			CreatedAt: dto.CreatedAtDTO{
				Seconds: hit.Source.CreatedAt.Seconds,
				Nanos:   hit.Source.CreatedAt.Nanos,
			},
		}
		products = append(products, product)
	}

	return products, nil
}
