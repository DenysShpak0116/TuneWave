package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/repository"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/digcontainer"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// @title API Specification
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide the Bearer token in the format: "Bearer {token}"
func main() {
	container := digcontainer.BuildContainer()
	// invoke container db
	err := container.Invoke(func(db *gorm.DB) {
		if err := repository.AutoMigrate(db); err != nil {
			log.Fatalf("Migration failed: %s", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to invoke DB migration: %s", err)
	}

	//invoke http server
	err = container.Invoke(func(router *chi.Mux, cfg *config.Config) {
		log.Println("Starting server on port", cfg.Http.Port)

		srv := &http.Server{
			Addr:    "localhost:" + strconv.Itoa(cfg.Http.Port),
			Handler: router,
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Listen: %s\n", err)
			}
		}()

		log.Println("Server started on port", cfg.Http.Port)

		<-quit
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
		log.Println("Server exited gracefully")
	})
	if err != nil {
		log.Fatalf("Failed to invoke HTTP server: %s", err)
	}
}
