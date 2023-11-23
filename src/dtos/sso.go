package dtos

type SsoLoginPayload struct {
	Client   string
	Username string
	Password string
}

type SsoLoginResponse struct {
	Callback string
	Exchange string
}
