package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/InspectorGadget/dockwatch-core/docker"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		// Fetch containers
		containers, err := docker.FetchContainers()
		if err != nil {
			log.Println("Error fetching containers:", err)
			break
		}

		// Send containers to the client
		err = conn.WriteJSON(containers)
		if err != nil {
			log.Println("Error writing to websocket:", err)
			break
		}

		time.Sleep(2 + time.Second)
	}
}
