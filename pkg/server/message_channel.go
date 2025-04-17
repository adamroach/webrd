package server

type MessageChannel interface {
	Send(message any) error
	Receive() (any, error)
	Close() error
}
