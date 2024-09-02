package grpc

import (
	"context"
	"log"
	"strconv"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/k1borgG/test_task/internal/dto"
	"github.com/k1borgG/test_task/internal/service"
	"github.com/k1borgG/test_task_grpc"
)

type GRPCServer struct {
	test_task_grpc.UnimplementedProductServiceServer
	productService service.ProductService
}

func NewGRPCServer(productService service.ProductService) *GRPCServer {
	return &GRPCServer{
		productService: productService,
	}
}

func (g *GRPCServer) AddProduct(ctx context.Context, req *test_task_grpc.AddProductRequest) (*test_task_grpc.AddProductResponse, error) {
	dtoReq := dto.MapAddProductRequestToDTO(req)
	dtoResp, err := g.productService.AddProduct(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	resp := &test_task_grpc.AddProductResponse{
		Id:      dtoResp.ID,
		Product: req.Product,
	}
	return resp, nil
}

func (g *GRPCServer) GetProduct(ctx context.Context, req *test_task_grpc.GetProductRequest) (*test_task_grpc.GetProductResponse, error) {
	dtoReq := dto.MapGetProductRequestToDTO(req)
	dtoResp, err := g.productService.GetProduct(ctx, dtoReq)
	if err != nil {
		log.Printf("error getting response: %s", err)
		return nil, err
	}

	var grpcProducts []*test_task_grpc.AddProduct
	for _, productDTO := range dtoResp.Products {
		price, _ := strconv.ParseUint(productDTO.Price, 10, 32)

		grpcProduct := &test_task_grpc.AddProduct{
			Name:        productDTO.Name,
			Description: productDTO.Description,
			Brand:       productDTO.Brand,
			Model:       productDTO.Model,
			Price:       uint32(price),
			Coordinates: &test_task_grpc.Coordinates{
				Lat: productDTO.Coordinates.Lat,
				Lon: productDTO.Coordinates.Lon,
			},
			CreatedAt: &timestamp.Timestamp{
				Seconds: productDTO.CreatedAt.Seconds,
				Nanos:   productDTO.CreatedAt.Nanos,
			},
		}
		grpcProducts = append(grpcProducts, grpcProduct)
	}

	resp := &test_task_grpc.GetProductResponse{
		Product: grpcProducts,
	}
	return resp, nil
}
