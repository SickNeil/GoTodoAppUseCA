// usecases/login_usecase.go
package usecases

import (
	"go-todo-app/entities"
	"go-todo-app/utils"
)

// UserProcessUseCase 定義登入相關的業務邏輯
type UserProcessUseCase struct {
	Repo entities.LoginRepository
}

// NewLoginUseCase 建立新的 UserProcessUseCase
func NewLoginUseCase(repo entities.LoginRepository) *UserProcessUseCase {
	return &UserProcessUseCase{Repo: repo}
}

// RegisterUser 用於註冊新使用者
func (uc *UserProcessUseCase) RegisterUser(username, password, email string) (string, error) {
	// 加密密碼
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := entities.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
	}
	return uc.Repo.CreateUser(user)
}

// Login 用於登入驗證
func (uc *UserProcessUseCase) Login(username, password string) (bool, error) {
	_, err := uc.Repo.GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	return false, nil
}
