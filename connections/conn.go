package connections

import (
    "github.com/go-redis/redis"
    "sync"
    "strconv"
    "errors"
    "github.com/redis-manager/go-redis-server/config"
    "github.com/jingweno/conf"
)



type RedisConnection struct {
    *redis.Client
    Conf *config.RedisConfig 	`json:"conf"`
    AllKeys [][]string	`json:"all_keys"`
    Err error			`json:"err"`
}

var onceRedis sync.Once
var redisCon *RedisConnection

func getRedisAndInitWithConfig(conf *config.RedisConfig) (*RedisConnection) {
    redisCon = &RedisConnection{
        Conf:conf,
    }

    redisCon.Client = redis.NewClient(&redis.Options{
        Addr:    conf.Host,
        Password: conf.Pw, // no password set
        DB:      conf.Db,  // use default DB
    })

    i := 0
    //try again
    for ;i<3; i++ {
        _, err := redisCon.Ping().Result()
        if err != nil {
            redisCon.reConnection()
            redisCon.Err = err
        }else{
            redisCon.Err = nil
        }
    }

    return redisCon
}

//初始化不同配置的连接池
func Init() map[string]*RedisConnection {
    onceRedis.Do(func() {
        connections := make(RedisConnections)
        for _,conf := range config.GetRedisHosts().GetAllServer() {
            hval := conf.GetHval()
            connections[hval] = GetRedis(hval)
            conn := connections[hval]
            if conn.Err == nil {
                conn.initKeys()
            }
        }
    })

    return connections
}

func GetRedisInConnections(hval string) *RedisConnection {
    conn,ok := connections[hval]
    if ok {
        return conn
    }

    connections[hval] = GetRedis(hval)

    return connections[hval]
}

func GetRedis(hval string)  *RedisConnection{
    conn,ok := connections[hval]
    if ok {
        return conn
    }

    conf,err := config.GetRedisHosts().GetConfig(hval)
    if  err != nil {
        return &RedisConnection{
            Err:errors.New("config not found with hval:"+hval),
        }
    }

    redisCon = &RedisConnection{
        Conf:conf,
    }

    redisCon.Client = redis.NewClient(&redis.Options{
        Addr:    conf.Host,
        Password: conf.Pw, // no password set
        DB:      conf.Db,  // use default DB
    })

    i := 0
    //try again
    for ;i<3; i++ {
        _, err1 := redisCon.Ping().Result()
        if err1 != nil {
            redisCon.reConnection()
            redisCon.Err = err1
        }else{
            redisCon.Err = nil
        }
    }
    connections[conf.GetHval()] = redisCon
    return redisCon
}

func (conn *RedisConnection) reConnection() {
    conf := conn.Conf
    redisCon.Client = redis.NewClient(&redis.Options{
        Addr:    conf.Host,
        Password: conf.Pw, // no password set
        DB:      conf.Db,  // use default DB
    })
}

//select dbs and keys to mem
func (conn *RedisConnection) initKeys() {
    s := conn.ConfigGet("databases")

    val,err := s.Result()
    db := 0
    if err == nil {
        if len(val) > 0 {
            db,err = strconv.Atoi(val[1].(string))
            if err != nil {
                db = 0
            }
        }
    }

    conn.AllKeys = make([][]string,0)

    for i:=0; i<=db; i++ {
        conn.Do("select", i)
        s := conn.Keys("*")
        keys,err := s.Result()
        if err != nil {
            conn.Err = err
        }else {
            conn.AllKeys[i] = keys
        }
    }
}

func RemoveConnections(hval string) error {
    err := connections.Remove(GetRedis(hval))

    if err != nil {
        return err
    }

    err = config.GetRedisHosts().Remove(hval)

    if err != nil {
        return err
    }


    return nil
}

func AddConnections(conf *config.RedisConfig) error {
    err := config.GetRedisHosts().Add(conf)

    if err != nil {
        return err
    }

    client := GetRedis(conf.GetHval())
    client.initKeys()

    if client.Err != nil {
        return client.Err
    }

    return nil
}

func ChangeConnection(hval string, conf *config.RedisConfig) error {
    err := RemoveConnections(hval)
    if err != nil {
        return err
    }

    return AddConnections(conf)
}