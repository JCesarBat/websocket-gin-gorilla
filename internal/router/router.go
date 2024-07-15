package router

import (
	"github.com/gin-gonic/gin"
	"webSocketGorrilaMuxGrpc/internal/user"
	"webSocketGorrilaMuxGrpc/internal/ws"
)

var r *gin.Engine

func InitRouter(handler *user.Handler, wsHandler *ws.HandlerWebSocket) {
	r = gin.Default()
	r.POST("/signup", handler.CreateUser)
	r.GET("/login", handler.Login)
	r.GET("/logout", handler.Logout)
	r.POST("/ws/CreateRoom", wsHandler.CreateRoom)
	r.GET("/ws/JoinRoom/:roomId", wsHandler.JoinRoom)

	r.GET("/ws/GetRooms", wsHandler.GetRooms)
	r.GET("/ws/GetClients/:roomId", wsHandler.GetClients)
}
func Start(addr string) error {
	return r.Run(addr)
}
