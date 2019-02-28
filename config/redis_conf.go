package config

import (
    "crypto/sha256"
    "time"
)

type RedisConfig struct {
    Name string `json:"name,omitempty"`
    Host string `json:"host,omitempty"`
    Db   int    `json:"db,omitempty"`
    Pw   string `json:"pw,omitempty"`
    Hval string `json:"hval,omitempty"`
}


func NewConf(host, name, pw string)*RedisConfig {
    return initConfig(host, name, pw)
}
func initConfig(host, name, pw string) *RedisConfig {
    conf := &RedisConfig{
        Name:name,
        Host:host,
        Db:0,
        Pw:pw,
    }
    conf.GetHval()
    return conf
}

func (RConf *RedisConfig) GetHval() string {
    if RConf.Hval == "" {
        bs := sha256.Sum256([]byte(RConf.Host + time.Now().String()))
        RConf.Hval = string(bs[:])
    }

    return RConf.Hval
}

func (RConf *RedisConfig) GetName() string {
    if RConf.Name == "" {
        return RConf.Host
    }
    return RConf.Name
}

