// entities/todo.go
package entities

import "time"

// Todo 定義業務的核心實體，可以用於任何資料庫
type Todo struct {
	ID        string    `json:"id" bson:"_id,omitempty"` // bson 標籤讓 MongoDB 能夠識別這個字段，但仍保持 string 類型
	Task      string    `json:"task"`
	CreatedAt time.Time `json:"created_at"`
}
