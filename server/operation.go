package server

import (
    "net/url"
    "strconv"
    "github.com/redis-manager/go-redis-server/connections"
    "github.com/redis-manager/go-redis-server/config"
    "github.com/redis-manager/go-redis-server/structure"
)

const(
    //config
    SetConfig = "set_config"
    ChangeConfig = "change_config"
    DelConfig =  "del_config"

    //server
    Info		= "get_server_info"

    //key
    GetKey    = "get_key"
    DelKey    = "del_key"
    SetTtl    = "set_ttl"

    //string
    SetValue  = "set_value"

    //set use it
    AddListValue  = "add_list_value"
    DelListValue  = "del_list_value"

    SetSortSetField  = "set_sort_set_field"
    DelSortSetField  = "del_sort_set_field"

    SetHashField  = "set_hash_field"
    DelHashField  = "del_hash_field"
)

func getServerInfo(conn *connections.RedisConnection) string {
    cmd := conn.Info()
    str,_ := cmd.Result()
    return str
}

func removeConfig(values url.Values) error {
    hval := values.Get("hval")

    return connections.RemoveConnections(hval)
}

func setConfig(values url.Values) error {
    host := values.Get("host")
    name := values.Get("name")
    pw   := values.Get("pw")

    return  connections.AddConnections(&config.RedisConfig{
        Name:name,
        Host:host,
        Db:0,
        Pw:pw,
    })
}

func changeConfig(values url.Values) error {
    hval := values.Get("hval")
    name := values.Get("name")
    pw   := values.Get("pw")
    host := values.Get("host")

    return connections.ChangeConnection(hval, config.NewConf(host, name, pw))
}

func KeyOperation(conn *connections.RedisConnection, action string, values url.Values) (kValues *structure.KeyValues) {
    dbstr := values.Get("db")
    db := 0
    if dbstr != "" {
        db,_ = strconv.Atoi(dbstr)
    }

    conn.Do("select", db)

    key := values.Get("key_name")
    info := &structure.KeyInfo{
        KeyName:key,
        Db:db,
        Type:getType(conn, key),
        TTl:getTtl(conn, key),
    }
    switch action {
    case GetKey:
        return getKey(conn, info)
    case DelKey:
    case SetTtl:
    case SetValue:
    case AddListValue:
    case DelListValue:
    case SetSortSetField:
    case DelSortSetField:
    case SetHashField:
    case DelHashField:
    default:
        return nil
    }
    return nil
}

func DoOperation(values url.Values) (bool, interface{}) {
    redisHash := values.Get("redis_hash")
    conf,err := config.GetRedisHosts().GetConfig(redisHash)
    if err != nil {
        return false,nil
    }

    conn 	:= connections.GetRedisInConnections(conf.GetHval())
    action 	:= values.Get("action")

    switch action {
    case SetConfig:
        return true,setConfig(values)
    case DelConfig:
        return true,removeConfig(values)
    case ChangeConfig:
        return true,changeConfig(values)
    case Info:
        return true,getServerInfo(conn)
    default:
        result := KeyOperation(conn, action, values)
        if result != nil {
            return true, result
        }
        return false, nil
    }

    return false, nil
}

func getType(conn *connections.RedisConnection, key string) uint8  {
    cmd := conn.Type(key)
    str,err := cmd.Result()

    if err != nil {
        return structure.RedisUnknow
    }

    return structure.GetTypeWithString(str)
}

func getTtl(conn *connections.RedisConnection, key string) int64{
    cmd := conn.TTL(key)
    d,err := cmd.Result()

    if err != nil {
        return -1
    }

    return d.Nanoseconds()
}

func getKey(con *connections.RedisConnection, info *structure.KeyInfo) *structure.KeyValues {
    cmd := con.Get(info.GetKeyName())
    str,_ := cmd.Result()

    return &structure.KeyValues{
        KeyInfo:info,
        Value: str,
    }
}




