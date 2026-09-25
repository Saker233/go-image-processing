package api

import (
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	db "github.com/Saker233/go-image-processing/internal/database"
	"github.com/Saker233/go-image-processing/internal/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store *db.Store
}

type userRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
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

func (s *Server) register(c *gin.Context) {
	var req userRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashed, err := util.CreateHashed([]byte(req.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateUserParams{
		Username:     req.Username,
		PasswordHash: string(hashed),
	}

	user, err := s.store.CreateUser(c, args)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	rsp := userResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusCreated, rsp)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
