// usecases/login_usecase.go
package usecases

import (
	"go-todo-app/entities"
	"go-todo-app/utils"
)

// LoginUseCase 定義登入相關的業務邏輯
type LoginUseCase struct {
	Repo entities.LoginRepository
}

// NewLoginUseCase 建立新的 LoginUseCase
func NewLoginUseCase(repo entities.LoginRepository) *LoginUseCase {
	return &LoginUseCase{Repo: repo}
}

// RegisterUser 用於註冊新使用者
func (uc *LoginUseCase) RegisterUser(username, password, email string) error {
	// 加密密碼
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := entities.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
	}
	return uc.Repo.CreateUser(user)
}

// Login 用於登入驗證
func (uc *LoginUseCase) Login(username, password string) (bool, error) {
	user, err := uc.Repo.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	// 驗證密碼
	if utils.CheckPasswordHash(password, user.Password) {
		return true, nil
	}
	return false, nil
}
