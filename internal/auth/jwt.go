package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour

	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateAccessToken(userId uint, email, name, role string) (string, error)
	GenerateRefreshToken(userId uint, email, name, role string) (string, error)
	ValidateToken(tokenStr string) (*JWTClaims, error)
}

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) JWTService {
	if secretKey == "" {
		secretKey = "default_secret_change_me"
	}
	return &jwtService{secretKey: secretKey}
}

func (js *jwtService) generateToken(userId uint, email, name, role, tokenType string, duration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:    userId,
		Name:      name,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "spotsync",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(js.secretKey))
}

func (js *jwtService) GenerateAccessToken(userId uint, email, name, role string) (string, error) {
	return js.generateToken(userId, email, name, role, TokenTypeAccess, accessTokenDuration)
}

func (js *jwtService) GenerateRefreshToken(userId uint, email, name, role string) (string, error) {
	return js.generateToken(userId, email, name, role, TokenTypeRefresh, refreshTokenDuration)
}

func (js *jwtService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(js.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
