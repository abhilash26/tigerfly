package main

import (
	"log"

	"github.com/abhilash26/tigerfly/internal/env"
	"github.com/abhilash26/tigerfly/internal/router"
	"github.com/abhilash26/tigerfly/internal/server"
	"github.com/abhilash26/tigerfly/internal/view"
	"github.com/abhilash26/tigerfly/routes"
)

func main() {
	// Load Environment Variables
	envPath := ".env"
	if err := env.LoadEnvFile(envPath); err != nil {
		log.Fatalf("Error loading %s file %v", envPath, err)
	}

	// Load Templates
	if err := view.PreloadAllTemplates(); err != nil {
		log.Fatalf("Error preloading templates: %v", err)
	}

	// Routing
	r := router.New()
	router.Static(r, "/static", "./static")

	routes.RegisterWeb(r)
	routes.RegisterAPI(r)

	s := server.New()
	s.Start(r)
}
