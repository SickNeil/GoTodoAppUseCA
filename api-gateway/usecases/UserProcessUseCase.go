// usecases/login_usecase.go
package usecases

import (
	"go-todo-app/entities"
	"net/http"
)

// UserProcessUseCase 定義登入相關的業務邏輯
type UserProcessUseCase struct {
	Repo entities.ILoginRepo
}

// NewLoginUseCase 建立新的 UserProcessUseCase
func NewLoginUseCase(repo entities.ILoginRepo) *UserProcessUseCase {
	return &UserProcessUseCase{Repo: repo}
}

// RegisterUser 用於註冊新使用者
func (uc *UserProcessUseCase) RegisterUser(username, password, email string) (*http.Response, error) {
	return uc.Repo.Register(username, password, email)
}

// Login 用於登入驗證
func (uc *UserProcessUseCase) Login(username, password string) (*http.Response, error) {
	return uc.Repo.Login(username, password)
}
