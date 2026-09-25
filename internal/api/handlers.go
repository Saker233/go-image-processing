package api

import (
	"net/http"
	"os"

	db "github.com/Saker233/go-image-processing/internal/database"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store *db.Store
}

func SetupServer(store *db.Store) {
	server := &Server{
		store: store,
	}

	r := gin.Default()

	r.GET("/health", server.getHealth)
	r.POST("/register", server.register)

	r.Run(os.Getenv("PORT"))
}

func (s *Server) getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Healthy",
	})
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
