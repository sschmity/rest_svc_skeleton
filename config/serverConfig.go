package config

import (
	"sync"
)

type ServerConfigSingleton struct {
	port string
}

const port=":13200"

var serverConfigInstance *ServerConfigSingleton
var onceServer sync.Once

func GetServerConfigInstance() *ServerConfigSingleton {
	onceServer.Do(func() {
		serverConfigInstance = &ServerConfigSingleton{port:port}
	})
	return serverConfigInstance
}

func (ServerConfigSingleton) GetServerPort() string {
	return serverConfigInstance.port
}
func (ServerConfigSingleton) SetServerPort(aPort string) {
	serverConfigInstance.port=aPort
}