package auth

type Secret interface {
	Get() any
}
