package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type HandlerWebSocket struct {
	hub *Hub
}

func NewHandlerWebSocket(h *Hub) *HandlerWebSocket {
	return &HandlerWebSocket{
		hub: h,
	}
}

type CreateRoomRequest struct {
	ID   string `json:"id"  binding:"required"`
	Name string `json:"name"  binding:"required"`
}

func (h *HandlerWebSocket) CreateRoom(ctx *gin.Context) {
	var req CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}
	ctx.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *HandlerWebSocket) JoinRoom(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	roomId := ctx.Param("roomId")
	ClientId := ctx.Query("userId")
	username := ctx.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       ClientId,
		RoomID:   roomId,
		Username: username,
	}
	m := &Message{
		Content:  "A new user has Join the room",
		RoomID:   roomId,
		Username: username,
	}
	h.hub.Register <- cl

	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *HandlerWebSocket) GetRooms(c *gin.Context) {
	rooms := make([]RoomResponse, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomResponse{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

type GetClientRequest struct {
	RoomId string `uri:"roomId" binding:"required"`
}
type ClientResponse struct {
	ID       string `json:"ID"`
	Username string `json:"Username"`
}

func (h *HandlerWebSocket) GetClients(c *gin.Context) {
	var id GetClientRequest
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, ok := h.hub.Rooms[id.RoomId]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "the room does not exist"})
	}
	clients := []ClientResponse{}

	for _, c := range h.hub.Rooms[id.RoomId].Clients {
		clients = append(clients, ClientResponse{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
