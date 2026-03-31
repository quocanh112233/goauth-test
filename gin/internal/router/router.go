package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocanh112233/goauth-test/gin/config"
	"github.com/quocanh112233/goauth-test/gin/internal/handler"
	"github.com/quocanh112233/goauth-test/gin/internal/middleware"
	"github.com/quocanh112233/goauth-test/gin/internal/renderer"
	"github.com/quocanh112233/goauth-test/gin/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(cfg *config.Config, db *mongo.Database, authService service.AuthService) *gin.Engine {
	r := gin.Default()

	r.HTMLRender = renderer.NewHTMLRenderer(cfg.TemplateDir)

	authHandler := handler.NewAuthHandler(authService, cfg)
	homeHandler := handler.NewHomeHandler(cfg)
	dashboardHandler := handler.NewDashboardHandler(cfg)
	apiHandler := handler.NewAPIHandler()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	r.GET("/health", func(c *gin.Context) {
		status := "ok"
		dbStatus := "connected"
		if err := db.Client().Ping(context.Background(), nil); err != nil {
			status = "error"
			dbStatus = "disconnected"
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    status,
				"framework": "Gin",
				"db":        dbStatus,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    status,
			"framework": "Gin",
			"db":        dbStatus,
		})
	})

	r.GET("/login", authHandler.ShowLogin)
	r.POST("/login", authHandler.Login)
	r.GET("/signup", authHandler.ShowSignup)
	r.POST("/signup", authHandler.Signup)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(authService, cfg))
	{
		protected.GET("/home", homeHandler.ShowHome)
		protected.GET("/dashboard", dashboardHandler.ShowDashboard)
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/api/me", apiHandler.GetMe)
	}

	return r
}
