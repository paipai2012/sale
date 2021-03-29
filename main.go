package main

import (
	"context"
	"fmt"
	"log"
	"moose-go/engine"
	"moose-go/router"
	"moose-go/util"

	"github.com/gin-gonic/gin"

	socketio "github.com/googollee/go-socket.io"
)

var (
	ctx = context.Background()
)

func main() {
	initApp()
}

func initApp() {

	config, err := util.ParseConfig("./config/moose.json")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_, err = engine.NewOrmEngine(config)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	rdb := engine.NewRedisHelper()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal(err.Error())
		return
	}

	gin.SetMode(config.AppMode)

	app := gin.Default()

	// 使用中间件
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	// app.Use(Auth())

	router.InitRouter(app)

	// init socket io
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	go server.Serve()
	defer server.Close()
	router.InitSocket(app, server)

	app.Run("0.0.0.0:8090")
}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("auth ...")
	}
}
