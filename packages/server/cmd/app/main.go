package main

import (
	"log"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/repository"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/digcontainer"
	"gorm.io/gorm"
)

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
}
