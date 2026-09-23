package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func SetupServer() {
	godotenv.Load("app.env")
	r := gin.Default()
	r.GET("/health", getHealth)
	r.Run(os.Getenv("PORT"))
}

func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Healthy",
	})
}
