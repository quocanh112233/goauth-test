package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quocanh112233/goauth-test/gin/config"
	"github.com/quocanh112233/goauth-test/gin/internal/pkg/jwt"
	"github.com/quocanh112233/goauth-test/gin/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(authService service.AuthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err == nil {
			claims, err := jwt.ParseToken(accessToken, cfg.JWTSecret)
			if err == nil {
				userID, _ := primitive.ObjectIDFromHex(claims.UserID)
				user, err := authService.GetUserByID(c.Request.Context(), userID)
				if err == nil && user != nil {
					c.Set("user", user)
					c.Next()
					return
				}
			}
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err == nil {
			newAccessToken, err := authService.RefreshAccessToken(c.Request.Context(), refreshToken)
			if err == nil {
				c.SetCookie("access_token", newAccessToken, 1800, "/", "", cfg.IsProduction, true)

				claims, _ := jwt.ParseToken(newAccessToken, cfg.JWTSecret)
				userID, _ := primitive.ObjectIDFromHex(claims.UserID)
				user, _ := authService.GetUserByID(c.Request.Context(), userID)
				if user != nil {
					c.Set("user", user)
					c.Next()
					return
				}
			}
		}

		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			c.Redirect(http.StatusSeeOther, "/login")
		}
		c.Abort()
	}
}
