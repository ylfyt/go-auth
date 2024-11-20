package shared

type EnvConf struct {
	Port                      int
	DbConnection              string `json:"-"`
	JwtAccessTokenExpiryTime  int    // in seconds
	JwtAccessTokenSecretKey   string `json:"-"`
	JwtRefreshTokenExpiryTime int    // in minutes
	JwtRefreshTokenSecretKey  string `json:"-"`
}
