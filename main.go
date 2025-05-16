package main

import (
	"fmt"
	"os"

	"github.com/InspectorGadget/dockwatch-core/controllers"
	"github.com/InspectorGadget/dockwatch-core/docker"
	"github.com/gin-gonic/gin"
)

func init() {
	err := docker.ConnectToDockerSock()
	if err != nil {
		fmt.Println(fmt.Errorf("an error has occurred: %v", err.Error()))
		os.Exit(1)
	}

	defer docker.GetClient().Close()
}

func main() {
	r := gin.Default()

	r.GET("/socket", controllers.WSHandler)

	if err := r.Run(":8080"); err != nil {
		fmt.Printf("error has occurred: %v", err.Error())
	}
}
