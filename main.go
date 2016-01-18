package main

import (
	"log"
	"net/http"

	"time"

	"github.com/appwilldev/Instafig/conf"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
)

func main() {
	if conf.IsEasyDeployMode() && !conf.IsMasterNode() {
		if err := slaveCheckMaster(); err != nil {
			log.Panicf("slave node failed to check master: %s", err.Error())
		}

		go func() {
			for {
				time.Sleep(60 * time.Second)
				slaveCheckMaster()
			}
		}()
	}

	if conf.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginIns := gin.New()
	ginIns.Use(gin.Recovery())
	if conf.DebugMode {
		ginIns.Use(gin.Logger())
	}

	// client api
	clientAPIGroup := ginIns.Group("/client")
	{
		clientAPIGroup.GET("/conf", ClientReqData)
	}

	// op api
	opAPIGroup := ginIns.Group("/op")
	{
		opAPIGroup.GET("/users/:page", GetUsers)
		opAPIGroup.POST("/user", ConfWriteCheck, NewUser)

		opAPIGroup.GET("/apps/:user_key", GetApps)
		opAPIGroup.POST("/app", ConfWriteCheck, NewApp)
		opAPIGroup.PUT("/app", ConfWriteCheck, UpdateApp)

		opAPIGroup.GET("/configs/:app_key", GetConfigs)
		opAPIGroup.POST("/config", ConfWriteCheck, NewConfig)
		opAPIGroup.PUT("/config", ConfWriteCheck, UpdateConfig)
	}

	if conf.IsEasyDeployMode() {
		ginInsNode := gin.New()
		if conf.DebugMode {
			ginInsNode.Use(gin.Logger())
		}
		ginInsNode.Use(gin.Recovery())
		ginInsNode.POST("/node/req/:req_type", NodeRequestHandler)

		gracehttp.Serve(
			&http.Server{Addr: conf.HttpAddr, Handler: ginIns},
			&http.Server{Addr: conf.NodeAddr, Handler: ginInsNode})
	} else {
		gracehttp.Serve(&http.Server{Addr: conf.HttpAddr, Handler: ginIns})
	}
}
