package server


import (
    "net/http"
    "fmt"
    "text/template"
    "encoding/json"
    "github.com/redis-manager/go-redis-server/structure"
    "github.com/redis-manager/go-redis-server/config"
)


type Message struct {
    Url string
    index string
    Root  string
    FileHandler http.Handler
}

type Render struct {
    Key string
    Value structure.ValueOf
}



func (message *Message) ServeHTTP(res http.ResponseWriter,req *http.Request) {

    //if static file go file or some query
    if req.Method == "GET" {
        //default index
        if req.RequestURI == "" {
            req.RequestURI = "/index.html"
        }

        //show all keys
        if req.RequestURI == "/all" {
            _, err := template.ParseFiles(message.Root+"/resources/app/index.html")
            if err != nil {
                fmt.Println("parse file err:", err)
                return
            }


            res.WriteHeader(200)

            //for _,conf :=  range config.GetRedisHosts().GetAllServer() {
                //r := GetRedis(conf.GetHval())
                //r.initKeys()


                //if err := t.Execute(res, re); err != nil {
                //	res.Write([]byte(err.Error()))
                //	fmt.Println("There was an error:", err.Error())
                //}
                //return
            //}
        }

        //show config router
        if req.RequestURI == "/config" {
            configs := config.GetRedisHosts().GetAllServer()
            fmt.Println(configs)
            b,_ := json.Marshal(configs)
            res.Write(b)
            res.WriteHeader(200)
            return
        }
        message.FileHandler.ServeHTTP(res, req)
        return
    }

    //any request must use POST
    if req.Method != "POST" {
        //do any thing
        req.ParseForm()
        operation,result := DoOperation(req.PostForm)
        if !operation {
            res.Write([]byte("<h1>404</h1>"))
            res.WriteHeader(404)
        }

        bs,err := json.Marshal(result)

        if err != nil {
            res.Write([]byte(err.Error()))
            res.WriteHeader(500)
            return
        }

        res.Write(bs)
        res.WriteHeader(200)
        return
    }


    res.Write([]byte("<h1>404</h1>"))
    res.WriteHeader(404)

}

func (message *Message) getIndexContent()  {
    if message.index == "" {

    }else{

    }
}
