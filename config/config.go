package config

import "sync"

type HostConfig struct {
    Host string `json:"host,omitempty"`
    Port int    `json:"port,omitempty"`
}

var hostConfig *HostConfig
var hostOnce sync.Once
// default port 9987
func GetHostConfig() *HostConfig{
    hostOnce.Do(func(){
        hostConfig =  &HostConfig{
            Host:"127.0.0.1",
            Port:9987,
        }
    })
   return hostConfig
}

