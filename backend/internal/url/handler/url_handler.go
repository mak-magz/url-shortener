package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mak-magz/url-shortener/internal/url/model"
	"github.com/mak-magz/url-shortener/internal/url/service"

	appErrors "github.com/mak-magz/url-shortener/platform/errors"
)

type Handler interface {
	CreateShortURL(c *gin.Context)
	RedirectToOriginalURL(c *gin.Context)
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
		c.Error(appErrors.FormatValidationError(err))
		return
	}

	url, err := u.service.CreateShortURL(c, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": url})
}

// RedirectToOriginalURL implements [Handler].
func (u *URLHandler) RedirectToOriginalURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	originalURL, err := u.service.GetOriginalURL(c, shortCode)
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
