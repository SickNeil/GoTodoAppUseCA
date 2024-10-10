// interfaces/login/postgres_login_repository.go
package interfaces

import (
	"database/sql"
	"go-todo-app/entities"
)

type PostgresLoginRepository struct {
	DB *sql.DB
}

// NewPostgresLoginRepository 建立新的 PostgreSQL 登入儲存庫
func NewPostgresLoginRepository(db *sql.DB) *PostgresLoginRepository {
	return &PostgresLoginRepository{DB: db}
}

// CreateUser 用於在資料庫中創建新使用者
func (repo *PostgresLoginRepository) CreateUser(user entities.User) error {
	_, err := repo.DB.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)",
		user.Username, user.Password, user.Email)
	return err
}

// GetUserByUsername 根據使用者名稱取得使用者資訊
func (repo *PostgresLoginRepository) GetUserByUsername(username string) (entities.User, error) {
	var user entities.User
	err := repo.DB.QueryRow("SELECT id, username, password, email FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	return user, err
}

// VerifyUser 驗證使用者名稱和密碼是否正確
func (repo *PostgresLoginRepository) VerifyUser(username, password string) (bool, error) {
	user, err := repo.GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	if user.Password == password {
		return true, nil
	}
	return false, nil
}
