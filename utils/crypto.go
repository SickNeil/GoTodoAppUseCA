// utils/crypto.go
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 接受一個純文字密碼，返回加密後的字串
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost 為預設的工作因數，值為 10，數值越高越安全，但計算時間也越長
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 比較純文字密碼和加密後的密碼是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
