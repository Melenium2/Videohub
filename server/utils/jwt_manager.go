package utils

import (
	"VideoHub/database/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// JWTManager jwt token manager
type JWTManager struct {
	secret        string
	tokenDuration time.Duration
}

// UserClaims contains some user information
type UserClaims struct {
	jwt.StandardClaims
	ID       int64
	Username string
	Role     string
}

// NewJWTManager create new instance of jwt token manager
func NewJWTManager(secret string, duration time.Duration) *JWTManager {
	return &JWTManager{
		secret:        secret,
		tokenDuration: duration,
	}
}

// Verify check is token valid
func (manager *JWTManager) Verify(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errorUnexpectedJwtSigningMethod()
		}

		return []byte(manager.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errorInvalidTokenClaims()
	}

	return claims, nil
}

// Sign seals the token with user claims
func (manager *JWTManager) Sign(user *models.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		ID: user.ID,
		Username: user.Username,
		Role: user.Role,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(manager.secret))
}
