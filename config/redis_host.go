package config

import (
    "sync"
    "errors"
)

type RedisHostMap map[string]*RedisConfig

type IRedisHosts interface {
    GetConfig(hash string) (*RedisConfig,error)
    Add(*RedisConfig) error
    Remove(hash string) error
    Edit(hash string, New *RedisConfig) error
    GetAllServer()  map[string]*RedisConfig
}

var redisHosts RedisHostMap
var redisHostOnce sync.Once

func GetRedisHosts() IRedisHosts {
    redisHostOnce.Do(func() {
        if redisHosts == nil {
            redisHosts = make(RedisHostMap)
        }

    })

    return redisHosts
}


func (host RedisHostMap) GetConfig(hash string) (*RedisConfig, error) {
    conf, ok := host[hash]
    if ok {
        return conf, nil
    }

    return nil, errors.New("config not found with " + hash)
}

func (host RedisHostMap) GetConfigByName(name string) (*RedisConfig, error) {
    for _, conf := range host {
        if conf.GetName() == name {
            return conf, nil
        }
    }

    return nil, errors.New("config not found with " + name)
}

func (host RedisHostMap) GetName(hash string) (string, error) {
    conf, err := host.GetConfig(hash)
    if err != nil {
        return "", err
    }

    return conf.GetName(), nil
}

func (host RedisHostMap) Add(RConf *RedisConfig) (error) {
    for _, conf := range host {
        if conf.GetHval() == RConf.GetHval() {
            return errors.New("config already exists with name(if name is nil, name is host:port):" + RConf.GetName())
        }
    }

    host[RConf.GetHval()] = RConf
    return nil
}

func (host RedisHostMap) Remove(hash string) (error) {
    deleted := false

    for key, conf := range host {
        if conf.GetHval() == hash {
            delete(host, key)
            deleted = true
            break
        }
    }

    if !deleted {
        return errors.New("delete config error with hash:" + hash)
    }

    return nil
}

func (host RedisHostMap) Edit(hash string, New *RedisConfig) (error) {
    changed := false

    for key, conf := range host {
        if conf.GetHval() == hash {
            delete(host, key)
            host[New.GetHval()] = New
            changed = true
            break
        }
    }

    if !changed {
        return errors.New("change config error with hval:" + hash)
    }

    return nil
}

func (host RedisHostMap) GetAllServer() map[string]*RedisConfig {
    return host
}