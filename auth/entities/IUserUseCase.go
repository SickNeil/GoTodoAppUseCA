package entities

type IUserUseCase interface {
	Login(username, password string) (string, error)
	Register(user *User) error
	IsTokenValid(token string) (bool, error)
}
