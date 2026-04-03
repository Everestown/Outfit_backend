package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return false
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}
