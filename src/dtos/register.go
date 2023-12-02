package dtos

type Register struct {
	Username string `validate:"min=4,max=8"`
	Password string `validate:"min=4,max=8"`
}
