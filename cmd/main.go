package main

import (
	"log"

	"github.com/arjunsaxaena/Moonrider-Assignment/controllers"
	"github.com/arjunsaxaena/Moonrider-Assignment/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	repository.InitDB()

	r := gin.Default()
	r.POST("/identify", controllers.IdentifyContact)

	log.Println("Server running on port 8080")
	r.Run(":8080")
}
