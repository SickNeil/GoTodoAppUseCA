package entities

// 可以將 JWT token 設置到對象中的接口
type JwtSettable interface {
	SetJWT(token string)
}
