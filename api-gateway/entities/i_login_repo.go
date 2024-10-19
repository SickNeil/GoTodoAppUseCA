// entities/i_login_repo.go
package entities

import "net/http"

// LoginRepository 定義登入相關的資料庫操作介面
type LoginRepository interface {
	CreateUser(user User) (string, error)
	GetUserByUsername(username string) (UserInfo, error)
	VerifyUser(token string) (bool, error)
	Register(userName, password, email string) (*http.Response, error)
	Login(userName, password string) (*http.Response, error)
}
