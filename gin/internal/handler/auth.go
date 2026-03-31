package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocanh112233/goauth-test/gin/config"
	"github.com/quocanh112233/goauth-test/gin/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	if _, exists := c.Get("user"); exists {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{"Framework": h.cfg.Framework})
}

func (h *AuthHandler) ShowSignup(c *gin.Context) {
	if _, exists := c.Get("user"); exists {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	c.HTML(http.StatusOK, "signup.html", gin.H{"Framework": h.cfg.Framework})
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	accessToken, refreshToken, role, err := h.authService.Login(c.Request.Context(), email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error":     err.Error(),
			"Framework": h.cfg.Framework,
		})
		return
	}

	c.SetCookie("access_token", accessToken, 1800, "/", "", h.cfg.IsProduction, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/", "", h.cfg.IsProduction, true)

	if role == "admin" {
		c.Redirect(http.StatusSeeOther, "/dashboard")
	} else {
		c.Redirect(http.StatusSeeOther, "/home")
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	confirm := c.PostForm("confirm_password")

	if password != confirm {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"Error":     "Passwords do not match",
			"Framework": h.cfg.Framework,
		})
		return
	}

	if err := h.authService.Signup(c.Request.Context(), name, email, phone, password); err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"Error":     err.Error(),
			"Framework": h.cfg.Framework,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		_ = h.authService.Logout(c.Request.Context(), refreshToken)
	}

	c.SetCookie("access_token", "", -1, "/", "", h.cfg.IsProduction, true)
	c.SetCookie("refresh_token", "", -1, "/", "", h.cfg.IsProduction, true)

	c.Redirect(http.StatusSeeOther, "/login")
}
