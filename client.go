package main

import (
	"github.com/appwilldev/Instafig/conf"
	"github.com/gin-gonic/gin"
)

func ClientReqData(c *gin.Context) {
	clientData := &ClientData{
		AppKey:     c.Query("app_key"),
		OSType:     c.Query("os_type"),
		OSVersion:  c.Query("os_version"),
		AppVersion: c.Query("app_version"),
		Ip:         c.Query("ip"),
		Lang:       c.Query("lang"),
		DeviceId:   c.Query("device_id"),
	}

	if conf.IsEasyDeployMode() {
		if !conf.IsMasterNode() {
			//todo: to check node sync status
		}

		confData := getAppMatchConf(clientData.AppKey, clientData)
		memConfMux.RLock()
		nodes := memConfNodes
		memConfMux.RUnlock()

		nodeRes := make([]string, len(nodes))
		ix := 0
		for _, node := range nodes {
			nodeRes[ix] = node.URL
			ix++
		}

		Success(c, map[string]interface{}{
			"nodes": nodeRes,
			"confs": confData,
		})
		return
	}
}
