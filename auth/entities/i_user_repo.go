package entities

type UserRepository interface {
	Login(username, password string) (string, error)
	Register(user *User) error
}
