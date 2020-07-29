package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()

	r.GET("/ping", pingAPI)

	r.GET("/zoom/meetings", getMeetings)
	r.POST("/zoom/meetings", createMeeting)

	r.GET("/ews/get-user-availability", getUserAvailability)

	r.Run(":" + os.Getenv("GIN_PORT"))
}

func pingAPI(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
