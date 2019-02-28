package go_redis_server

import (
    "github.com/redis-manager/go-redis-server/config"
    "net/url"
    "path/filepath"
    "os"
    "github.com/redis-manager/go-redis-server/server"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/redis-manager/go-redis-server/connections"
    "strconv"
)

var debug = true
func Run()  {
    config.GetRedisHosts()


    buildHttpServerHandler()
    //if err != nil {
    //	fmt.Println(err)
    //}
    //buildEctron(urlStr)
    //select {
    //
    //}
}

func loadPath(path string) {
    var redisConfigList =  make([]*config.RedisConfig,0)
    data, err := ioutil.ReadFile(path)
    if err != nil {
        data, _ = ioutil.ReadFile("config.json")
    }
    b := []byte(data)
    err = json.Unmarshal(b, redisConfigList)
    if err == nil {
        for _,conf := range redisConfigList {
            err := connections.AddConnections(conf)
            if err != nil {
                //todo: do any thing
            }
        }
    }
}

func buildHttpServerHandler() (err error) {
    serverConfig :=config.GetHostConfig()
    urlStr := "http://"+serverConfig.Host+":"+strconv.Itoa(serverConfig.Port)
    urlObj,err := url.Parse(urlStr)
    if err != nil {
        return err
    }

    var root string
    if !debug {
        root, err = filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            return err
        }
    }else{
        root = os.Getenv("GOPATH")+"/src/github.com/lwl1989/go-redis-manager"
    }
    loadPath(root+"/config/config.json")
    //fmt.Println(root)
    message := &server.Message{
        Url: urlStr,
        Root: root,
        FileHandler:http.FileServer(http.Dir(root+"/resources/app")),
    }
    //fmt.Println(urlObj.Port())
    //+urlObj.Port()
    server.Init()
    err = http.ListenAndServe(":"+urlObj.Port(), message)
    if err != nil {
        fmt.Println(err)
        return err
    }
    return nil
}