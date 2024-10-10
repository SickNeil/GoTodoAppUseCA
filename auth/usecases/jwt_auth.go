// usecases/jwt_auth.go
package usecases

import (
	"auth/entities"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth interface {
	GenerateToken(user *entities.User) (string, error)
	IsTokenValid(token string) (bool, error)
}

type jwtAuth struct {
	secret []byte
}

func NewJWTAuth() JWTAuth {
	secret := os.Getenv("JWT_SECRET")
	fmt.Println("secret", secret)
	if secret == "" {
		secret = "defaultKey"
	}
	return &jwtAuth{
		secret: []byte(secret),
	}
}

func (j *jwtAuth) GenerateToken(user *entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *jwtAuth) IsTokenValid(tokenString string) (bool, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "defaultKey"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, nil
}
