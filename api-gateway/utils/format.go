package utils

import "time"

// 格式化時間的函數
func FormatAsDate(t time.Time) string {
	// 創建一個 UTC+8 的時區
	location := time.FixedZone("UTC+8", 8*60*60)

	// 將時間轉換為 UTC+8
	utc8Time := t.In(location)

	// 格式化為您需要的格式
	return utc8Time.Format("2006-01-02 15:04:05")
}
