// entities/user.go
package entities

// User 定義使用者實體，用於表示資料庫中的 users 表
type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
