package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type CustomClaims struct {
	AccountID    string `json:"account_id"`
	UniversityID string `json:"university_id"`
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

// Tạo token
func GenerateToken(accountID, userID, universityID primitive.ObjectID, role string, duration time.Duration) (string, error) {
	claims := CustomClaims{
		AccountID:    accountID.Hex(),
		UserID:       userID.Hex(),
		UniversityID: universityID.Hex(), // thêm trường này
		Role:         role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Parse token và lấy claims
func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token không hợp lệ")
	}

	return claims, nil
}
