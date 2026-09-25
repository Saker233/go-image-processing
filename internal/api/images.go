package api

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"net/http"
	"strconv"
	"strings"

	db "github.com/Saker233/go-image-processing/internal/database"
	"github.com/Saker233/go-image-processing/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type imageResponse struct {
	ID               uuid.UUID      `json:"id"`
	UserID           uuid.UUID      `json:"user_id"`
	OriginalFilename string         `json:"original_filename"`
	OriginalS3Key    string         `json:"original_s3_key"`
	ProcessedS3Key   sql.NullString `json:"processed_s3_key,omitempty"`
	MimeType         string         `json:"mime_type"`
	OriginalSize     int64          `json:"original_size"`
	ProcessedSize    sql.NullInt64  `json:"processed_size,omitempty"`
	Width            int32          `json:"width"`
	Height           int32          `json:"height"`
	CreatedAt        string         `json:"created_at"`
}

type transformRequest struct {
	Transformations map[string]interface{} `json:"transformations"`
}

func (s *Server) createImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authenticated")))
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid user identity")))
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer uploadedFile.Close()

	contents := &bytes.Buffer{}
	if _, err := contents.ReadFrom(uploadedFile); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	data := contents.Bytes()
	mimeType := http.DetectContentType(data)
	if !strings.HasPrefix(mimeType, "image/") {
		c.JSON(http.StatusUnsupportedMediaType, errorResponse(errors.New("file must be an image")))
		return
	}

	key, err := util.UploadS3File(bytes.NewReader(data), file.Filename, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	img, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("unable to read image dimensions: %w", err)))
		return
	}

	stored, err := s.store.CreateImage(c, db.CreateImageParams{
		UserID:           uid,
		OriginalFilename: file.Filename,
		OriginalS3Key:    key,
		MimeType:         mimeType,
		OriginalSize:     int64(len(data)),
		Width:            int32(img.Width),
		Height:           int32(img.Height),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":                stored.ID,
		"user_id":           stored.UserID,
		"original_filename": stored.OriginalFilename,
		"original_s3_key":   stored.OriginalS3Key,
		"mime_type":         stored.MimeType,
		"original_size":     stored.OriginalSize,
		"width":             stored.Width,
		"height":            stored.Height,
		"created_at":        stored.CreatedAt,
	})
}

func (s *Server) listImages(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authenticated")))
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid user identity")))
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	images, err := s.store.GetImagesByUser(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	start := (page - 1) * limit
	if start >= len(images) {
		c.JSON(http.StatusOK, gin.H{"images": []db.Image{}, "page": page, "limit": limit, "total": len(images)})
		return
	}

	end := start + limit
	if end > len(images) {
		end = len(images)
	}

	c.JSON(http.StatusOK, gin.H{
		"images": images[start:end],
		"page":   page,
		"limit":  limit,
		"total":  len(images),
	})
}

func (s *Server) getImage(c *gin.Context) {
	imageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authenticated")))
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid user identity")))
		return
	}

	img, err := s.store.GetImageByID(c, imageID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if img.UserID != uid {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("you do not own this image")))
		return
	}

	c.JSON(http.StatusOK, img)
}

func (s *Server) deleteImage(c *gin.Context) {
	imageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authenticated")))
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid user identity")))
		return
	}

	img, err := s.store.GetImageByID(c, imageID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if img.UserID != uid {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("you do not own this image")))
		return
	}

	if err := util.DeleteS3Object(img.OriginalS3Key); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if img.ProcessedS3Key.Valid {
		if err := util.DeleteS3Object(img.ProcessedS3Key.String); err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	if err := s.store.DeleteImage(c, imageID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image deleted"})
}

func (s *Server) transformImage(c *gin.Context) {
	imageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req transformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authenticated")))
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid user identity")))
		return
	}

	img, err := s.store.GetImageByID(c, imageID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if img.UserID != uid {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("you do not own this image")))
		return
	}

	if len(req.Transformations) == 0 {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("at least one transformation is required")))
		return
	}

	updated, err := s.store.UpdateProcessedImage(c, db.UpdateProcessedImageParams{
		ID:             img.ID,
		ProcessedS3Key: sql.NullString{String: img.OriginalS3Key, Valid: true},
		ProcessedSize:  sql.NullInt64{Int64: img.OriginalSize, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "image transformed",
		"image":           updated,
		"transformations": req.Transformations,
	})
}
