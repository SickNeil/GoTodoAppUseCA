package entities

// ITodoRepo 是一個用於處理待辦事項資料的接口
type ITodoRepo interface {
	Insert(todo Todo) error
	GetAll() ([]Todo, error)
	Delete(id string) error
}
