// usecases/user_usecase.go
package usecases

import (
	"auth/entities"
	"auth/interfaces"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo interfaces.PostgresUserRepo
	jwtAuth  JWTAuth
}

func NewUserUseCase(userRepo interfaces.PostgresUserRepo, jwtAuth JWTAuth) entities.IUserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

func (u *UserUseCase) Login(username, password string) (string, error) {
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

func (u *UserUseCase) Register(user *entities.User) error {
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

func (u *UserUseCase) IsTokenValid(token string) (bool, error) {
	return u.jwtAuth.IsTokenValid(token)
}
