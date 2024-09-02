package dto

import (
	"strconv"

	"github.com/k1borgG/test_task_grpc"
)

func MapAddProductRequestToDTO(req *test_task_grpc.AddProductRequest) AddProductRequestDTO {
	return AddProductRequestDTO{
		Product: ProductDTO{
			Name:        req.Product.GetName(),
			Description: req.Product.GetDescription(),
			Brand:       req.Product.GetBrand(),
			Model:       req.Product.GetModel(),
			Price:       strconv.Itoa(int(req.Product.GetPrice())),
			Coordinates: CoordinatesDTO{
				Lat: req.Product.GetCoordinates().GetLat(),
				Lon: req.Product.GetCoordinates().GetLon(),
			},
			CreatedAt: CreatedAtDTO{
				Seconds: req.Product.CreatedAt.Seconds,
				Nanos:   req.Product.CreatedAt.Nanos,
			},
		},
	}
}

func MapGetProductRequestToDTO(req *test_task_grpc.GetProductRequest) GetProductRequestDTO {
	var name, description, brand, model, filterPrice1, filterPrice2, filterDate1, filterDate2 *string

	// Проверка поля Name
	if req.Product.GetName() != "" {
		name = new(string)
		*name = req.Product.GetName()
	}

	// Проверка поля Description
	if req.Product.GetDescription() != "" {
		description = new(string)
		*description = req.Product.GetDescription()
	}

	// Проверка поля Brand
	if req.Product.GetBrand() != "" {
		brand = new(string)
		*brand = req.Product.GetBrand()
	}

	// Проверка поля Model
	if req.Product.GetModel() != "" {
		model = new(string)
		*model = req.Product.GetModel()
	}

	// Проверка поля Price
	var price *uint32
	if req.Product.GetPrice() != 0 {
		price = new(uint32)
		*price = req.Product.GetPrice()
	}

	// Проверка диапазонов фильтрации по цене
	if req.Product.GetFilterPrice1() != "" {
		filterPrice1 = new(string)
		*filterPrice1 = req.Product.GetFilterPrice1()
	}
	if req.Product.GetFilterPrice2() != "" {
		filterPrice2 = new(string)
		*filterPrice2 = req.Product.GetFilterPrice2()
	}

	// Проверка диапазонов фильтрации по дате
	if req.Product.GetFilterDate1() != "" {
		filterDate1 = new(string)
		*filterDate1 = req.Product.GetFilterDate1()
	}
	if req.Product.GetFilterDate2() != "" {
		filterDate2 = new(string)
		*filterDate2 = req.Product.GetFilterDate2()
	}

	// Инициализация CreatedAtDTO с проверкой на nil
	var createdAtDTO CreatedAtDTO
	if req.Product.GetCreatedAt() != nil {
		createdAtDTO = CreatedAtDTO{
			Seconds: req.Product.GetCreatedAt().Seconds,
			Nanos:   req.Product.GetCreatedAt().Nanos,
		}
	}

	// Инициализация CreatedAtDTO с проверкой на nil
	var coordinatesDTO CoordinatesDTO
	if req.Product.GetCoordinates() != nil {
		coordinatesDTO = CoordinatesDTO{
			Lat: req.Product.GetCoordinates().Lat,
			Lon: req.Product.GetCoordinates().Lon,
		}
	}

	// Создание DTO
	return GetProductRequestDTO{
		Name:         name,
		Description:  description,
		Brand:        brand,
		Model:        model,
		Price:        price,
		FilterPrice1: filterPrice1,
		FilterPrice2: filterPrice2,
		FilterDate1:  filterDate1,
		FilterDate2:  filterDate2,
		Coordinates:  coordinatesDTO,
		CreatedAt:    createdAtDTO,
	}
}
