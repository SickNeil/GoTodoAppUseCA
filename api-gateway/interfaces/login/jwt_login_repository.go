package interfaces

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-todo-app/entities"
	"net/http"
	"os"
	"time"
)

// 實作 JwtLoginRepository 介面
type JwtLoginRepository struct {
	Client *http.Client
}

// NewJwtLoginRepository 用於建立 JwtLoginRepository 實例
func NewJwtLoginRepository() *JwtLoginRepository {
	return &JwtLoginRepository{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// CreateUser 用於新增使用者
func (repo *JwtLoginRepository) CreateUser(user entities.User) (string, error) {
	// 送出request到auth-service
	jsonData, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	authServerUrl := os.Getenv("AUTH_SERVER_URL")
	resp, err := repo.Client.Post(authServerUrl+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 檢查回應狀態碼是否為 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("登入失敗: " + resp.Status)
	}

	// 解析回應，取得 JWT
	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	token, ok := result["token"]
	if !ok {
		return "", errors.New("無法取得 JWT Token")
	}

	return token, nil
}

// GetUserByUsername 用於取得使用者資訊
func (repo *JwtLoginRepository) GetUserByUsername(username string) (entities.UserInfo, error) {
	authServerUrl := os.Getenv("AUTH_SERVER_URL")
	// 送出request到auth-service
	resp, err := repo.Client.Get(authServerUrl + "/user/" + username)
	if err != nil {
		return entities.UserInfo{}, err
	}
	defer resp.Body.Close()

	// 檢查回應狀態碼是否為 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return entities.UserInfo{}, errors.New("無法取得使用者資訊: " + resp.Status)
	}

	// 解析回應，取得使用者資訊
	var user entities.UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return entities.UserInfo{}, err
	}

	return user, nil
}

// VerifyUser 用於驗證使用者
func (repo *JwtLoginRepository) VerifyUser(token string) (bool, error) {
	// 送出request到auth-service
	jwtData := map[string]string{
		"token": token,
	}
	jsonData, err := json.Marshal(jwtData)
	if err != nil {
		return false, err
	}

	authServerUrl := os.Getenv("AUTH_SERVER_URL")
	resp, err := repo.Client.Post(authServerUrl+"/validate-token", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 檢查回應狀態碼是否為 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return false, errors.New("登入失敗: " + resp.Status)
	}

	// 解析回應，取得驗證結果
	var result map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result["result"], nil
}
