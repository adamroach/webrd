package auth

type Authenticator interface {
	Authenticate(username, password string) (token string, err error)
	ValidateToken(token string) (username string, err error)
}
