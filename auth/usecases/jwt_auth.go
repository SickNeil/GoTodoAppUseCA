package usecases

import (
	"auth/entities"
	"fmt"
	"os"
	"time"

	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth interface {
	GenerateToken(user *entities.User) (string, error)
	IsTokenValid(token string) (bool, error)
}

type jwtAuth struct {
	privateKey *rsa.PrivateKey
}

// IsTokenValid implements JWTAuth.
func (j *jwtAuth) IsTokenValid(tokenString string) (bool, error) {
	// 讀取公鑰檔案
	keyData, err := os.ReadFile("/key/public.key")
	if err != nil {
		return false, fmt.Errorf("failed to read public key: %v", err)
	}

	// 解析 Token 並驗證簽名
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 確認 token 使用的演算法是 RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// 解析 RSA 公鑰
		return jwt.ParseRSAPublicKeyFromPEM(keyData)
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return false, err
	}

	// 檢查 token 是否有效
	if !token.Valid {
		fmt.Println("Invalid token")
		return false, fmt.Errorf("token is invalid")
	}

	return true, nil
}

func NewJWTAuth() JWTAuth {
	// 加載私鑰
	keyData, err := os.ReadFile("/key/secret.key")
	if err != nil {
		panic("Error reading private key: " + err.Error())
	}

	// 解析私鑰
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		panic("Error parsing private key: " + err.Error())
	}

	return &jwtAuth{
		privateKey: privateKey,
	}
}

func (j *jwtAuth) GenerateToken(user *entities.User) (string, error) {
	// 使用 RS256 來生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// 使用私鑰簽名
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
