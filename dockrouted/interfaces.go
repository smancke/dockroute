package dockrouted

//go:generate go get github.com/golang/mock/mockgen
//go:generate mockgen -self_package dockrouted -package dockrouted -destination interfaces_mocks_test.go github.com/smancke/dockroute/dockrouted Backend
//go:generate sed -ie "s/dockrouted .dockroute\\/dockrouted.//g;s/dockrouted\\.//g" interfaces_mocks_test.go
type Backend interface {
	GetService(name string) (hostAndPort string, err error)
}
