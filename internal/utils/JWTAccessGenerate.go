package utils

import (
	"context"
	"encoding/base64"
	"go-auth-service/config"
	"go-auth-service/models/token"
	"strings"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTAccessClaims jwt claims
type JWTAccessClaims struct {
	token.AccessTokenClaims
}

// Valid claims verification
func (a *JWTAccessClaims) Valid() error {
	if time.Unix(a.ExpireTime, 0).Before(time.Now()) {
		return errors.ErrInvalidAccessToken
	}
	return nil
}

func NewSM2JWTAccessGenerate(config *config.AppConfig, method jwt.SigningMethod) *JWTAccessGenerate {
	return &JWTAccessGenerate{
		config:       config,
		SignedKeyID:  "",
		SignedKey:    nil,
		SignedMethod: method,
	}
}

// NewJWTAccessGenerate create to generate the jwt access token instance
func NewJWTAccessGenerate(kid string, key []byte, method jwt.SigningMethod) *JWTAccessGenerate {
	return &JWTAccessGenerate{
		config:       nil,
		SignedKeyID:  kid,
		SignedKey:    key,
		SignedMethod: method,
	}
}

// JWTAccessGenerate generate the jwt access token
type JWTAccessGenerate struct {
	config       *config.AppConfig
	SignedKeyID  string
	SignedKey    []byte
	SignedMethod jwt.SigningMethod
}

// Token based on the UUID generated token
func (a *JWTAccessGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	claims := &token.AccessTokenClaims{}
	claims.Issuer = data.UserID
	claims.ID = uuid.New().String()
	claims.ExpireTime = data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix()
	claims.ExpiresAt = jwt.NewNumericDate(data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()))
	claims.IssuedAt = jwt.NewNumericDate(data.TokenInfo.GetAccessCreateAt())
	claims.NotBefore = jwt.NewNumericDate(data.TokenInfo.GetAccessCreateAt())
	claims.Subject = data.Client.GetID()
	claims.Scope = data.TokenInfo.GetScope()
	claims.Audience = jwt.ClaimStrings{data.Client.GetDomain()}
	accessToken := jwt.NewWithClaims(a.SignedMethod, claims)
	if a.SignedKeyID != "" {
		accessToken.Header["kid"] = a.SignedKeyID
	}
	var key interface{}
	if a.isEs() {
		v, err := jwt.ParseECPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", err
		}
		key = v
	} else if a.isRsOrPS() {
		v, err := jwt.ParseRSAPrivateKeyFromPEM(a.SignedKey)
		if err != nil {
			return "", "", err
		}
		key = v
	} else if a.isHs() {
		key = a.SignedKey
	} else if a.isSM2() {
		v := LoadSM2PrivateKey(a.config)
		key = v
	} else {
		return "", "", errors.New("unsupported sign method")
	}

	access, err := accessToken.SignedString(key)
	if err != nil {
		return "", "", err
	}
	refresh := ""

	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}

func (a *JWTAccessGenerate) isEs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "ES")
}

func (a *JWTAccessGenerate) isSM2() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "SM2")
}

func (a *JWTAccessGenerate) isRsOrPS() bool {
	isRs := strings.HasPrefix(a.SignedMethod.Alg(), "RS")
	isPs := strings.HasPrefix(a.SignedMethod.Alg(), "PS")
	return isRs || isPs
}

func (a *JWTAccessGenerate) isHs() bool {
	return strings.HasPrefix(a.SignedMethod.Alg(), "HS")
}
