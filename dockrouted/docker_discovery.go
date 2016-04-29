package dockrouted

type DockerDiscovery struct {
}

func NewDockerDiscovery() *DockerDiscovery {
	return &DockerDiscovery{}
}

func (d *DockerDiscovery) GetService(name string) (host string, port string, err error) {
	return "www.google.de", "80", nil
}
