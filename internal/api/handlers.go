package api

import (
	"errors"
	"net/http"
	"os"
	"strings"

	db "github.com/Saker233/go-image-processing/internal/database"
	"github.com/Saker233/go-image-processing/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	r.POST("/login", server.login)

	auth := r.Group("/")
	auth.Use(requireAuth())
	auth.POST("/images", server.createImage)
	auth.GET("/images", server.listImages)
	auth.GET("/images/:id", server.getImage)
	auth.DELETE("/images/:id", server.deleteImage)
	auth.POST("/images/:id/transform", server.transformImage)

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

func requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization header required")))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("use format: Bearer <token>")))
			return
		}

		claims, err := util.ParseJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid user identity in token")))
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
