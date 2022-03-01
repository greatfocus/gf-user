package server

import (
	"net/http"
	"strings"
	"time"

	jwt5 "github.com/golang-jwt/jwt/v5"
)

// JWT
type JWT interface {
	CreateToken(tokenInfo TokenInfo) (string, error)
	IsValidToken(r *http.Request) bool
	GetTokenInfo(r *http.Request) (*TokenInfo, error)
	Secret() string
}

// TokenInfo struct
type TokenInfo struct {
	Permissions []string
	Origin      string
	ActorID     int64
}

// jwt struct
type jwt struct {
	authorized bool
	minutes    int64
	secret     string
}

// NewJWT
func NewJWT(secret string, minutes int64, authorized bool) JWT {
	return &jwt{
		minutes:    minutes,
		authorized: authorized,
		secret:     secret,
	}
}

// CreateToken generates jwt for API login
func (j *jwt) CreateToken(tokenInfo TokenInfo) (string, error) {
	claims := &jwt5.MapClaims{
		"iss":  "issuer",
		"exp":  time.Now().Add(time.Minute * time.Duration(j.minutes)).Unix(),
		"data": tokenInfo,
	}
	token := jwt5.NewWithClaims(jwt5.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// IsValidToken checks for jwt validity
func (j *jwt) IsValidToken(r *http.Request) bool {
	tokenString := getToken(r)
	token, err := jwt5.Parse(tokenString, func(token *jwt5.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return false
	}
	claims := token.Claims.(jwt5.MapClaims)
	data := claims["data"].(TokenInfo)
	return r.Header.Get("Origin") != data.Origin
}

// GetTokenInfo returns token information
func (j *jwt) GetTokenInfo(r *http.Request) (*TokenInfo, error) {
	tokenString := getToken(r)
	token, err := jwt5.Parse(tokenString, func(token *jwt5.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt5.MapClaims)
	data := claims["data"].(TokenInfo)
	return &data, nil
}

// GetSecret returns secret information
func (j *jwt) Secret() string {
	return j.secret
}

// getToken get jwt from header
func getToken(r *http.Request) string {
	TokenToken := r.Header.Get("Authorization")
	if len(strings.Split(TokenToken, " ")) == 2 {
		return strings.Split(TokenToken, " ")[1]
	}
	return ""
}
