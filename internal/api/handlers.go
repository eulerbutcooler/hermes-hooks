package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type ExecutionEvent struct {
	RelayID    string          `json:"relay_id"`
	Payload    json.RawMessage `json:"payload"`
	ReceivedAt time.Time       `json:"received_at"`
}

type EventProducer interface {
	Publish(relayID string, event ExecutionEvent) error
}

type Handler struct {
	producer EventProducer
}

func NewHandler(p EventProducer) *Handler {
	return &Handler{producer: p}
}

func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	relayID := chi.URLParam(r, "relayID")
	if relayID == "" {
		http.Error(w, "Relay ID is required", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	event := ExecutionEvent{
		RelayID:    relayID,
		Payload:    body,
		ReceivedAt: time.Now(),
	}
	if err := h.producer.Publish(relayID, event); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"queued"`))
}
