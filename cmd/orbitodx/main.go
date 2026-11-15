package main

import (
	"flag"
	"fmt"
	"log"
	"orbitodx/internal/handlers"
	"orbitodx/internal/store"
	"orbitodx/web"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	defaultPort := os.Getenv("PORT")
	if defaultPort == "" {
		defaultPort = "8080"
	}

	host := flag.String("host", "0.0.0.0", "Host interface to bind to")
	port := flag.String("port", defaultPort, "Port to listen on")
	flag.Parse()

	appState := store.NewAppState()
	h := handlers.NewServerHandlers(appState)

	app := fiber.New(fiber.Config{
		AppName: "OrbitoDX v3.4 — NASA HRP & CSA Autonomous Crew Health Triage",
	})

	// Serve embedded static assets (/static/css/app.css, /static/js/cockpit.js)
	app.Get("/static*", static.New("", static.Config{
		FS:     web.StaticFS(),
		MaxAge: 3600,
	}))

	// Cockpit UI & Telemetry API Routes
	app.Get("/", h.HandleDashboard)
	app.Get("/api/health", h.HandleHealth)
	app.Get("/api/state", h.HandleGetState)
	app.Post("/api/evaluate", h.HandleEvaluateVitals)
	app.Post("/api/scenario/:id", h.HandleScenario)
	app.Post("/api/mission/:id", h.HandleSwitchMission)
	app.Post("/api/protocol/:stepId/toggle", h.HandleToggleStep)

	addr := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("🚀 OrbitoDX Autonomous Crew Health Cockpit listening on http://%s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
