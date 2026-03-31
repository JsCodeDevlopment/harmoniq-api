package main

import (
	"api/src/app"
	"log"
	"os"
)

func main() {
	engine := app.Bootstrap()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Application is starting on port %s...", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
