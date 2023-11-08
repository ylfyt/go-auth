package dtos

type RefreshPayload struct {
	RefreshToken string
}
type AccessPayload struct {
	RefreshToken string
}

type TokenPayload struct {
	RefreshToken string
	AccessToken  string
	ExpiredIn    int64
}
