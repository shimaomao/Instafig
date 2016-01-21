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
		DataSign:   c.Query("data_sign"),
	}

	if conf.IsEasyDeployMode() {
		if !conf.IsMasterNode() {
			//todo: to check node sync status
		}

		memConfMux.RLock()
		nodes := make([]string, len(memConfNodes))
		ix := 0
		for _, node := range memConfNodes {
			nodes[ix] = node.URL
			ix++
		}
		// do not support app data_sign in server-side, always return app configs
		needConf := true || (memConfApps[clientData.AppKey] != nil && clientData.DataSign != memConfApps[clientData.AppKey].DataSign)
		memConfMux.RUnlock()

		if needConf {
			configs := getAppMatchConf(clientData.AppKey, clientData)
			memConfMux.RLock()
			dataSign := memConfApps[clientData.AppKey].DataSign
			memConfMux.RUnlock()

			Success(c, map[string]interface{}{
				"nodes":     nodes,
				"configs":   configs,
				"data_sign": dataSign,
			})
		} else {
			Success(c, map[string]interface{}{
				"nodes": nodes,
			})
		}
		return
	}
}
