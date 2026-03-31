package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocanh112233/goauth-test/gin/config"
	"github.com/quocanh112233/goauth-test/gin/internal/model"
)

type HomeHandler struct {
	cfg *config.Config
}

func NewHomeHandler(cfg *config.Config) *HomeHandler {
	return &HomeHandler{cfg: cfg}
}

func (h *HomeHandler) ShowHome(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"User":      user.(*model.User),
		"Framework": h.cfg.Framework,
	})
}
