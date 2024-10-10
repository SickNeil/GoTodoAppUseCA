// entities/i_login_repo.go
package entities

// LoginRepository 定義登入相關的資料庫操作介面
type LoginRepository interface {
	CreateUser(user User) (string, error)
	GetUserByUsername(username string) (UserInfo, error)
	VerifyUser(token string) (bool, error)
}
