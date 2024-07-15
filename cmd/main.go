package main

import (
	"log"
	db "webSocketGorrilaMuxGrpc/db/sqlc"
	"webSocketGorrilaMuxGrpc/internal/router"
	"webSocketGorrilaMuxGrpc/internal/user"
	"webSocketGorrilaMuxGrpc/internal/ws"
)

func main() {

	store, err := db.NewStore()
	if err != nil {
		log.Fatalf("cannot connecto to the database %s", err)
	}
	handler := user.NewHandler(store)

	hub := ws.NewHub()
	wsHandler := ws.NewHandlerWebSocket(hub)
	router.InitRouter(handler, wsHandler)

	go hub.Run()

	err = router.Start("0.0.0.0:3000")
	if err != nil {
		panic(err)
	}
}
