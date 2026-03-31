package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocanh112233/goauth-test/gin/internal/model"
)

type APIHandler struct{}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

func (h *APIHandler) GetMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user.(*model.User))
}
