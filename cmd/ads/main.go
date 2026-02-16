package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	//include handlers and configuration
	"ADS4/internal/app"
	"ADS4/internal/config"
	"ADS4/internal/utils"
)

func main() {

	log.Printf("Starting ADS4 service")

	// Load the configuration
	cfg := config.LoadConfig()

	// Initialize the app
	application := app.NewApp(cfg)

	// Get PORT from environment, default to 8080 if not set
	port := cfg.ADSPORT
	if port == "" {
		port = "8080"
	}

	if len(os.Args) > 1 {
		s := os.Args[1]

		if _, err := strconv.ParseInt(s, 10, 64); err == nil {
			log.Printf("Using port %s", s)
			port = s
		}
	}

	// get the local IP that has Internet connectivity
	ip := utils.GetLocalIP()
	if ip == nil {
		ip = net.ParseIP("0.0.0.0")
	}
	log.Printf("Starting HTTP service on http://%s:%s", ip, port)
	//log.Printf("Shutdown the service http://%s:%s/shutdown (admin only)", ip, port)

	// HTTP listener is in a goroutine as it's blocking
	go func() {
		if err := application.Router.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err)

		}
	}()

	// Setup a ctrl-c trap to ensure a graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	application.Context = ctx

	defer cancel()

	log.Println("closing database connections")
	application.DB.Close()

	// Log the shutdown process
	log.Println("Shutting HTTP service down")
	if err := application.Router.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	log.Println("Shutdown complete")
	os.Exit(0)
}
