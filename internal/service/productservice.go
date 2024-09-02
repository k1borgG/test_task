package service

import (
	"context"
	"log"
	"strconv"

	"github.com/k1borgG/test_task/internal/dto"
	"github.com/k1borgG/test_task/internal/repository"
)

type ProductService struct {
	esRepo repository.ElasticsearchRepository
}

func NewProductService(esRepo repository.ElasticsearchRepository) ProductService {
	return ProductService{
		esRepo: esRepo,
	}
}

func (s *ProductService) AddProduct(ctx context.Context, req dto.AddProductRequestDTO) (*dto.AddProductResponseDTO, error) {
	productID, err := s.esRepo.IndexProduct(ctx, req.Product)
	if err != nil {
		return nil, err
	}

	return &dto.AddProductResponseDTO{
		ID:      productID,
		Product: req.Product,
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req dto.GetProductRequestDTO) (*dto.GetProductResponseDTO, error) {
	var mustQueries []map[string]interface{}
	var filterQueries []map[string]interface{}

	if req.Name != nil && *req.Name != "" {
		mustQueries = append(mustQueries, map[string]interface{}{
			"match": map[string]interface{}{"name": *req.Name},
		})
	}

	if req.Description != nil && *req.Description != "" {
		mustQueries = append(mustQueries, map[string]interface{}{
			"match": map[string]interface{}{"description": *req.Description},
		})
	}

	if req.Brand != nil && *req.Brand != "" {
		mustQueries = append(mustQueries, map[string]interface{}{
			"match": map[string]interface{}{"brand": *req.Brand},
		})
	}

	if req.Model != nil && *req.Model != "" {
		mustQueries = append(mustQueries, map[string]interface{}{
			"match": map[string]interface{}{"model": *req.Model},
		})
	}

	if req.Price != nil && *req.Price != 0 {
		mustQueries = append(mustQueries, map[string]interface{}{
			"match": map[string]interface{}{"price": *req.Price},
		})
	}

	if req.FilterPrice1 != nil && req.FilterPrice2 != nil {
		fPrice1, _ := strconv.Atoi(*req.FilterPrice1)
		fPrice2, _ := strconv.Atoi(*req.FilterPrice2)
		filterQueries = append(filterQueries, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte": fPrice1,
					"lte": fPrice2,
				},
			},
		})
	}

	if &req.Coordinates != nil {
		if req.Coordinates.Lat != 0 || req.Coordinates.Lon != 0 {
			mustQueries = append(mustQueries, map[string]interface{}{
				"geo_distance": map[string]interface{}{
					"distance": "200km",
					"coordinates": map[string]interface{}{
						"lat": req.Coordinates.Lat,
						"lon": req.Coordinates.Lon,
					},
				},
			})
		}
	}

	if req.FilterDate1 != nil && req.FilterDate2 != nil {
		fDate1, _ := strconv.Atoi(*req.FilterDate1)
		fDate2, _ := strconv.Atoi(*req.FilterDate2)
		filterQueries = append(filterQueries, map[string]interface{}{
			"range": map[string]interface{}{
				"created_at.seconds": map[string]interface{}{
					"gte": fDate1,
					"lte": fDate2,
				},
			},
		})
	}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   mustQueries,
				"filter": filterQueries,
			},
		},
	}

	products, err := s.esRepo.SearchProducts(ctx, searchQuery)
	if err != nil {
		log.Printf("Error searching products: %v", err)
		return nil, err
	}

	return &dto.GetProductResponseDTO{
		Products: products,
	}, nil
}
