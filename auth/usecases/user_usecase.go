// usecases/user_usecase.go
package usecases

import (
	"auth/entities"
	"auth/interfaces"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Login(username, password string) (string, error)
	Register(user *entities.User) error
}

type userUseCase struct {
	userRepo interfaces.PostgresUserRepository
	jwtAuth  JWTAuth
}

func NewUserUseCase(userRepo interfaces.PostgresUserRepository, jwtAuth JWTAuth) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

func (u *userUseCase) Login(username, password string) (string, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("使用者名稱或密碼錯誤")
	}

	// 使用 utils.ComparePassword 比較密碼
	exists, err := user.ComparePassword(password)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("使用者名稱或密碼錯誤")
	}

	// 簽發 JWT
	token, err := u.jwtAuth.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *userUseCase) Register(user *entities.User) error {
	// 檢查使用者是否已存在
	exists, err := u.userRepo.UserExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("使用者名稱已存在")
	}

	// 使用 utils.HashPassword 對密碼進行哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 創建使用者
	return u.userRepo.CreateUser(user)
}

func HashPassword(s string) {
	// 實作密碼哈希邏輯

}
