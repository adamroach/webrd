package auth

type PasswordChecker interface {
	CurrentUser() string
	CheckPassword(username, password string) bool
}
