package manager

import (
	"NewCsTeamServer/client"
	"github.com/gin-gonic/gin"
)

func handleGetClientList(c *gin.Context) {
	list := client.GlobalClientManager.GetClientList()
	//c.JSON(200, public.ApiResponse{
	//	Code:  0,
	//	Msg:   "ok",
	//	Count: int64(len(list)),
	//	Data:  list,
	//})
	c.HTML(200, "client_list.tmpl", gin.H{
		"Clients": list,
	})
	return
}
