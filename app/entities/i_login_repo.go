// entities/i_login_repo.go
package entities

// LoginRepository 定義登入相關的資料庫操作介面
type LoginRepository interface {
	CreateUser(user User) error
	GetUserByUsername(username string) (User, error)
	VerifyUser(username, password string) (bool, error)
}
