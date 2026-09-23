package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupServer() {
	r := gin.Default()
	r.GET("/health", getHealth)
	r.Run(":8000")
}

func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Healthy",
	})
}
