// services/auth_service.go
package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthService interface {
	Login(username, password string) (string, error)
}

type authService struct {
	AuthServerURL string
}

func NewAuthService(authServerURL string) AuthService {
	return &authService{AuthServerURL: authServerURL}
}

func (s *authService) Login(username, password string) (string, error) {
	// 構建登入請求
	loginData := map[string]string{
		"username": username,
		"password": password,
	}
	jsonData, _ := json.Marshal(loginData)

	// 發送請求到 auth-server
	resp, err := http.Post(s.AuthServerURL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 處理回應
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("登入失敗")
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	token, ok := result["token"]
	if !ok {
		return "", errors.New("未取得 Token")
	}
	return token, nil
}
