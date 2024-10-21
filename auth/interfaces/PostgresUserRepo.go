package interfaces

import (
	"auth/entities"
	"database/sql"
)

type PostgresUserRepo struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		DB: db,
	}
}

func (r *PostgresUserRepo) GetUserByUsername(username string) (*entities.User, error) {
	user := &entities.User{}
	err := r.DB.QueryRow("SELECT id, username, password, email FROM users WHERE username = $1",
		username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepo) CreateUser(user *entities.User) error {
	_, err := r.DB.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)",
		user.Username, user.Password, user.Email)
	return err
}

func (r *PostgresUserRepo) UserExists(username string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow("SELECT exists(SELECT 1 FROM users WHERE username=$1)", username).
		Scan(&exists)
	return exists, err
}
