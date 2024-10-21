// entities/i_login_repo.go
package entities

import "net/http"

// ILoginRepo 定義登入相關的資料庫操作介面
type ILoginRepo interface {
	Register(userName, password, email string) (*http.Response, error)
	Login(userName, password string) (*http.Response, error)
}
