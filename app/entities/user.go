// entities/user.go
package entities

// User 定義使用者實體，用於表示資料庫中的 users 表
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
