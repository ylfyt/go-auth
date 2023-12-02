package shared

type EnvConf struct {
	ListenPort                int
	DbConnection              string
	JwtAccessTokenExpiryTime  int // in seconds
	JwtAccessTokenSecretKey   string
	JwtRefreshTokenExpiryTime int // in minutes
	JwtRefreshTokenSecretKey  string
}
