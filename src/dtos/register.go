package dtos

type Register struct {
	Username string `validate:"required,min=4,max=8"`
	Password string `validate:"required,min=4,max=8"`
}
