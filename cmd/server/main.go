package main

import (
	"log"
	"net/http"

	"github.com/eulerbutcooler/hermes-hooks/internal/api"
	"github.com/eulerbutcooler/hermes-hooks/internal/config"
	"github.com/eulerbutcooler/hermes-hooks/internal/queue"
)

func main() {
	cfg := config.LoadConfig()
	natsQueue, err := queue.NewNatsQueue(cfg.NatsUrl)
	if err != nil {
		log.Fatalf("Fatal: Could not connect to NATS: %v", err)
	}
	log.Println("Success: Connected to NATS JetStream")

	handler := api.NewHandler(natsQueue)
	r := api.NewRouter(handler)

	log.Printf("Starting Webhook Server on port: %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
