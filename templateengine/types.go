package templateengine

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	*jwt.RegisteredClaims
	Username string `json:"username"`
}

type RenderData struct {
	Claims  *Claims
	BaseURL string
	Data    any
}
