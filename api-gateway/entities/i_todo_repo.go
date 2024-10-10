package entities

// TodoRepository 是一個用於處理待辦事項資料的接口
type TodoRepository interface {
	Insert(todo Todo) error
	GetAll() ([]Todo, error)
	Delete(id string) error
}

// 可以將 JWT token 設置到對象中的接口
type JwtSettable interface {
	SetJWT(token string)
}
