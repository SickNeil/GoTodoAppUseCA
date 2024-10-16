// usecases/todo_usecase.go
package usecases

import (
	"fmt"
	"go-todo-app/entities"
	"time"
)

// TodoUseCase 是業務邏輯層的結構體
type TodoUseCase struct {
	Repo entities.TodoRepository // 使用接口而非具體實作
}

// NewTodoUseCase 創建新的 UseCase
func NewTodoUseCase(repo entities.TodoRepository) *TodoUseCase {
	return &TodoUseCase{Repo: repo}
}

// AddTodo 用於新增待辦事項
func (uc *TodoUseCase) AddTodo(task string) error {
	todo := entities.Todo{
		Task:      task,
		CreatedAt: time.Now(),
	}

	return uc.Repo.Insert(todo)
}

// GetTodos 用於取得所有待辦事項
func (uc *TodoUseCase) GetTodos() ([]entities.Todo, error) {
	return uc.Repo.GetAll()
}

// DeleteTodo 用於刪除待辦事項
func (uc *TodoUseCase) DeleteTodo(id string) error {
	return uc.Repo.Delete(id)
}

func (uc *TodoUseCase) SetJWT(repo entities.TodoRepository, token string) {
	fmt.Println("todo usecase SetJWT")
	if jwtSettable, ok := repo.(entities.JwtSettable); ok {
		fmt.Println("todo usecase SetJWT ok")
		jwtSettable.SetJWT(token)
	}
	fmt.Println("todo usecase SetJWT end")
}
