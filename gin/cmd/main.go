package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quocanh112233/goauth-test/gin/config"
	"github.com/quocanh112233/goauth-test/gin/db"
	"github.com/quocanh112233/goauth-test/gin/internal/repository"
	"github.com/quocanh112233/goauth-test/gin/internal/router"
	"github.com/quocanh112233/goauth-test/gin/internal/service"
)

func main() {
	cfg := config.Load()

	client, err := db.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}
	defer db.Disconnect(client)

	database := client.Database(cfg.MongoDB)

	userRepo := repository.NewUserRepository(database)
	sessionRepo := repository.NewSessionRepository(database)
	authService := service.NewAuthService(userRepo, sessionRepo, cfg.JWTSecret)

	r := router.Setup(cfg, database, authService)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
