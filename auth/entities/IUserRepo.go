package entities

type UserRepository interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(user *User) error
	UserExists(username string) (bool, error)
	Login(username, password string) (string, error)
	Register(user *User) error
	IsTokenValid(token string) (bool, error)
}
