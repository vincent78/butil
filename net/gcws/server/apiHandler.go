package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func wsInfo(c *gin.Context) {
	info := map[string]interface{}{}

	list1 := make([]string, len(HubManager.loginServers))
	for _, server := range HubManager.loginServers {
		list1 = append(list1, fmt.Sprintf("%v - %v", server.RemoteIP, server.ID))
	}
	info["loginServers"] = map[string]interface{}{
		"count":  len(HubManager.loginServers),
		"detail": list1,
	}

	list2 := make([]string, len(HubManager.servers))
	for server, _ := range HubManager.servers {
		list2 = append(list2, server.RemoteIP)
	}

	info["anonymity"] = map[string]interface{}{
		"count":  len(HubManager.servers),
		"detail": list2,
	}
	c.JSON(http.StatusOK, info)
}
