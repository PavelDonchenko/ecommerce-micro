package models

import "github.com/golang-jwt/jwt/v4"

type TokenClaims struct {
	Sub    string   `json:"sub,omitempty"`
	Email  string   `json:"email,omitempty"`
	Jti    string   `json:"jti,omitempty"`
	Nbf    int64    `json:"nbf,omitempty"`
	Iat    int64    `json:"iat,omitempty"`
	Exp    int64    `json:"exp,omitempty"`
	Iss    string   `json:"iss,omitempty"`
	Claims []Claims `json:"claims,omitempty"`
	jwt.RegisteredClaims
}

type Claims struct {
	Type  string `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
}
