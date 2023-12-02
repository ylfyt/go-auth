package dtos

type SsoLoginPayload struct {
	Client   string `validate:"required"`
	Username string `validate:"required,min=4,max=8"`
	Password string `validate:"required,min=4,max=8"`
}

type SsoLoginResponse struct {
	Callback string
	Exchange int64
}
