package dtos

type CreateProduct struct {
	Name        string `json:"name" validate:"min=5,max=225"`
	Description string `json:"description" validate:"min=5,max=1000"`
	Price       int64  `json:"price" validate:"gte=1"`
}
