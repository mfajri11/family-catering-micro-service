package security

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mfajri11/family-catering-micro-service/auth-service/config"
)

// there are 3 types of token `access token`, `password token` and `refresh token`
// each of them must has a unique purpose

var (
	iss    string = config.Cfg.App.Name
	secret string = config.Cfg.App.Security.JWTKey()
)

type AuthClaim struct {
	jwt.RegisteredClaims
	Type tokenType `json:"type"`
}

func newAuthClaim() *AuthClaim {
	now := timeNow()
	return &AuthClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(config.Cfg.App.AccessTokenDuration) * time.Minute)),
		},
	}
}

func (s security) GenerateToken(tokenType string) (string, error) {

	claim := newAuthClaim()
	claim.Type = newTokenType(tokenType)
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	if claim.Type == RT {
		return base64.URLEncoding.EncodeToString([]byte(tokenStr)), nil
	}

	return base64.StdEncoding.EncodeToString([]byte(tokenStr)), nil

}

func (s security) ValidateToken(tokenStr string) (any, error) {
	tokenByte, err := base64.StdEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	tokenStrDecoded := string(tokenByte)
	claim := &AuthClaim{}
	token, err := jwt.ParseWithClaims(tokenStrDecoded, claim, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("security.Validate: error unexpected method algorithm header got %s", t.Header["alg"])
		}

		return []byte(secret), nil

	})

	if err != nil {
		return nil, fmt.Errorf("security.Validate: error validating token: %w", err) // TODO: wrap error
	}

	payload, ok := token.Claims.(*AuthClaim)
	if !ok {
		return nil, errors.New("security.Validate: invalid claim type") // TODO: wrap error
	}
	if payload.Issuer != iss {
		return nil, errors.New("security.Validate: unexpected token") // TODO: wrap error
	}

	return payload, nil

}
