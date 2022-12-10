package dtos

type RefreshPayload struct {
	RefreshToken string `json:"refresh_token"`
}
type AccessPayload struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenPayload struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
