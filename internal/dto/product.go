package dto

type CoordinatesDTO struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type CreatedAtDTO struct {
	Seconds int64 `json:"seconds"`
	Nanos   int32 `json:"nanos"`
}

type ProductDTO struct {
	Name        string         `json:"name"`        // Название продукта
	Description string         `json:"description"` // Описание продукта
	Brand       string         `json:"brand"`       // Бренд продукта
	Model       string         `json:"model"`       // Модель продукта
	Price       string         `json:"price"`       // Цена продукта
	Coordinates CoordinatesDTO `json:"coordinates"` // Координаты продукта
	CreatedAt   CreatedAtDTO   `json:"created_at"`  // Время создания продукта
}

type AddProductRequestDTO struct {
	Product ProductDTO `json:"product"`
}

type AddProductResponseDTO struct {
	ID      string     `json:"id"`      // ID нового продукта
	Product ProductDTO `json:"product"` // Добавленный продукт
}

type GetProductRequestDTO struct {
	Name         *string        `json:"name,omitempty"`
	Description  *string        `json:"description,omitempty"`
	Brand        *string        `json:"brand,omitempty"`
	Model        *string        `json:"model,omitempty"`
	Price        *uint32        `json:"price,omitempty"`
	FilterDate1  *string        `json:"filter_date1,omitempty"`
	FilterDate2  *string        `json:"filter_date2,omitempty"`
	FilterPrice1 *string        `json:"filter_price1,omitempty"`
	FilterPrice2 *string        `json:"filter_price2,omitempty"`
	Coordinates  CoordinatesDTO `json:"coordinates"` // Координаты продукта
	CreatedAt    CreatedAtDTO   `json:"created_at"`  // Время создания продукта
}

type GetProductResponseDTO struct {
	Products []ProductDTO `json:"products"` // Список найденных продуктов
}
