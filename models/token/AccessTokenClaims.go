package token

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	ClientID   string `json:"client_id"`
	Scope      string `json:"scope"`
	ExpireTime int64  `json:"expire_time"`
}
