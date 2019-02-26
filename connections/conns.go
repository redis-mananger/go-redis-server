package connections

import (
    "github.com/redis-manager/go-redis-server/config"
    "errors"
)

type RedisConnections map[string]*RedisConnection
var connections RedisConnections


type IRedisConnections interface {
    Add(*config.RedisConfig) error
    Remove(conn *RedisConnection) error
}


func (rcons RedisConnections) Remove(conn *RedisConnection) error {
    deleted := false
    for key,conn := range rcons {
        if conn.Conf.GetHval() == conn.Conf.GetHval() {
            delete(rcons, key)
            deleted = true
            break
        }
    }

    if !deleted {
        return errors.New("remove connections error:"+conn.Conf.GetName())
    }

    return nil
}

func (rcons RedisConnections)  Add(conf *config.RedisConfig) error {
    found := false
    for _,conn := range rcons {
        if conn.Conf.GetHval() == conf.GetHval() {
            found = true
            break
        }
    }

    if !found {
        return errors.New("add connections error(already exists):"+conf.GetName())
    }

    con := getRedisAndInitWithConfig(conf)
    con.initKeys()
    return nil
}

