package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocanh112233/goauth-test/gin/config"
	"github.com/quocanh112233/goauth-test/gin/internal/model"
)

type DashboardHandler struct {
	cfg *config.Config
}

func NewDashboardHandler(cfg *config.Config) *DashboardHandler {
	return &DashboardHandler{cfg: cfg}
}

func (h *DashboardHandler) ShowDashboard(c *gin.Context) {
	u, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	user := u.(*model.User)
	if user.Role != "admin" {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"User":      user,
		"Framework": h.cfg.Framework,
	})
}
