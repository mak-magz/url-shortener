package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mak-magz/url-shortener/internal/url/model"
	"github.com/mak-magz/url-shortener/internal/url/service"
)

type Handler interface {
	CreateShortURL(c *gin.Context)
	GetOriginalURL(c *gin.Context)
}

type URLHandler struct {
	service service.Service
}

func NewURLHandler(service service.Service) Handler {
	return &URLHandler{service: service}
}

// CreateShortURL implements [Handler].
func (u *URLHandler) CreateShortURL(c *gin.Context) {
	var req model.CreateURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := u.service.CreateShortURL(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": url})
}

// GetOriginalURL implements [Handler].
func (u *URLHandler) GetOriginalURL(c *gin.Context) {
	panic("unimplemented")
}
