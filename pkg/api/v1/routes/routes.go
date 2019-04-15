package routes

import jwt "github.com/dgrijalva/jwt-go"

type RespMessage struct {
	Status  int
	Message string
}

type LoginRespMessage struct {
	Status       int
	Message      string
	Token        string
	RefreshToken string
}

// Token configuration
var MySigningKey = []byte("AllYourBase")

type TokenType string

const (
	NormalTokenType  TokenType = "normal"
	RefreshTokenType TokenType = "refresh"
)

type MyCustomClaims struct {
	Type     TokenType `json:"type"`
	UserName string    `json:"username"`
	jwt.StandardClaims
}
