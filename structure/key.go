package structure

import "time"

//redis type
const RedisNone   = 0
const RedisString = 1
const RedisList   = 2
const RedisSet    = 3
const RedisZset   = 4
const RedisHash   = 5
const RedisGeo    = 6
const RedisUnknow = 9

//db key map lists
type KeyMap map[int][]*KeyInfo

//key
type KeyInfo struct {
    Db		int    `json:"db"`   //which db
    KeyName string `json:"key_name"`  //keyname
    Type 	uint8  `json:"type"`      //key type
    TTl		int64  `json:"ttl"`
    checkExists  bool
}

func GetKeyInfo() *KeyInfo  {
    return &KeyInfo{
        Db: 0,
        KeyName:"",
        Type:RedisUnknow,
        TTl: -1,
        checkExists:false,
    }
}

func GetKeyInfoWithBasic(key string, db int) *KeyInfo  {
    return &KeyInfo{
        Db: db,
        KeyName:key,
        Type:RedisUnknow, //default string
        TTl: -1,
        checkExists:true,
    }
}

func GetKeyMap() KeyMap  {
    //keys := make([]*KeyInfo,0)
    return make(KeyMap)
}

func (maps KeyMap) GetDbKeys(db int) []*KeyInfo {
    keyList,ok := maps[db]
    if ok {
        return keyList
    }

    return nil
}

func (maps KeyMap) GetKeyWith(db int, keyName string) *KeyInfo {

    list := maps.GetDbKeys(db)

    if list == nil {
        return nil
    }

    for _,value := range list {
        if value.KeyName == keyName {
            return value
        }
    }

    return  nil
}

func (maps KeyMap) String() string {
    str := ""

    for db,value := range maps {
        str += "db:"+string(db)+" <br/>"
        for _,key := range  value {
            str += key.String()
        }
    }

    return str
}

func (info *KeyInfo) GetDb() int {
    return info.Db
}

func (info *KeyInfo) GetKeyName() string {
    return info.KeyName
}

func (info *KeyInfo) GetType() uint8 {
    return info.Type
}

func (info *KeyInfo) GetTypeString() string {
    switch info.Type {
    case RedisHash:
        return "hash"
    case RedisList:
        return "list"
    case RedisSet:
        return "set"
    case RedisZset:
        return "zset"
    case RedisString:
        return "string"
    case RedisGeo:
        return "geo"
    case RedisNone:
        fallthrough  //fallthrough不会判断下一条case的expr结果是否为true。 就是没有break
        //但是如果几个条件都走一样的结果，使用 fallthrough串起来即可
    default:
        return "none"
    }
}

func (info *KeyInfo) GetTtl() int64 {
    return info.TTl
}

func (info *KeyInfo) SetTtlWithTime(t time.Duration) {
    info.TTl = int64(t.Seconds())
}

func (info *KeyInfo) SetTypeWithString(t string) {
    if t == "none" {
        info.Type = RedisNone
        info.checkExists = false
    } else {
        info.Type = GetTypeWithString(t)
        info.checkExists = true
    }
}

func (info *KeyInfo) String() string {
    return "name:"+info.KeyName + " type:" + info.GetTypeString()
}


func GetTypeWithString(t string)  uint8 {
    switch t {
    case "hash":
        return RedisHash
    case "list":
        return RedisList
    case "geo":
        return RedisGeo
    case "string":
        return RedisString
    case "set":
        return RedisSet
    case "zset":
        return RedisZset
    case "none":
        return  RedisNone
    default:
        return RedisUnknow
    }
}
