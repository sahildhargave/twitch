package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var appID, appCertificate string

func main() {

	appIDEnv, appIDExists := os.LookupEnv("9fc6de75e3a14ff7a77c16c7eb6bb767")
	appCertEnv, appCertExists := os.LookupEnv("b3fab562898646baad8a7a9e5923b0b0")

	if !appIDExists || !appCertExists {
		log.Fatal("FATAL ERROR: ENV not properly configured, check APP_ID and APP_CERTIFICATE")
	} else {
		appID = appIDEnv
		appCertificate = appCertEnv
	}

	api := gin.Default()

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api.Run(":8080")
}
