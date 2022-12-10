package dtos

type RefreshPayload struct {
	RefreshToken string `json:"refreshToken"`
}
type AccessPayload struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenPayload struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
	ExpiredIn    int64  `json:"expiredIn"`
}
