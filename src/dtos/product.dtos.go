package dtos

type CreateProduct struct {
	Name        string `validate:"min=5,max=225"`
	Description string `validate:"min=5,max=1000"`
	Price       int64  `validate:"gte=1"`
}
