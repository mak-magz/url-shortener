package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appErrors "github.com/mak-magz/url-shortener/platform/errors"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			var appErr *appErrors.AppError

			if errors.As(err, &appErr) {
				c.JSON(appErr.Code, appErrors.AppError{
					Code:    appErr.Code,
					Message: appErr.Message,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}
	}
}
