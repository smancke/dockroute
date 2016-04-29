package dockrouted

type Backend interface {
	GetService(name string) (host string, port string, err error)
}
