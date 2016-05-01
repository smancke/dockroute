package dockrouted

import (
	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"sync"
)

type DockerDiscovery struct {
	server struct {
		mutex sync.Mutex
		names []string
	}
}

func NewDockerDiscovery() *DockerDiscovery {
	dd := &DockerDiscovery{}
	dd.server.names = make([]string, 0, 10)
	go dd.monitorHost("unix:///var/run/docker.sock")
	return dd
}

func (dd *DockerDiscovery) addServer(server string) {
	dd.server.mutex.Lock()
	defer dd.server.mutex.Unlock()

	for _, val := range dd.server.names {
		if server == val {
			// already contained, nothing to do
			return
		}
	}
	log.WithFields(log.Fields{
		"name":   server,
		"action": "add",
	}).Infof("adding server %#v", server)

	dd.server.names = append(dd.server.names, server)
}

func (dd *DockerDiscovery) removeServer(server string) {
	dd.server.mutex.Lock()
	defer dd.server.mutex.Unlock()

	for i, val := range dd.server.names {
		if server == val {
			log.WithFields(log.Fields{
				"name":   server,
				"action": "remove",
			}).Infof("removing server %#v", server)
			dd.server.names = append(dd.server.names[:i], dd.server.names[i+1:]...)
			return
		}
	}
}

func (dd *DockerDiscovery) monitorHost(host string) {
	docker, _ := dockerclient.NewDockerClient(host, nil)

	docker.StartMonitorEvents(func(event *dockerclient.Event, ec chan error, args ...interface{}) {
		switch event.Status {
		case "start":
			log.Printf("Received event: %#v\n", *event)
			dd.addServer(event.Actor.Attributes["name"])
		case "stop":
			log.Printf("Received event: %#v\n", *event)
			dd.removeServer(event.Actor.Attributes["name"])
		}
	}, nil)

	containers, err := docker.ListContainers(false, false, "")
	if err != nil {
		log.Error(err)
	}
	for _, c := range containers {
		dd.addServer(c.Names[0])
		log.Println(c.Id, c.Names)
	}
}

func (dd *DockerDiscovery) GetService(name string) (hostAndPort string, err error) {
	return "www.google.de:80", nil
}
