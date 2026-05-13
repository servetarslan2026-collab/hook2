package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DeliveryUpdate is the real-time event pushed to clients.
type DeliveryUpdate struct {
	Type           string    `json:"type"`
	ID             uuid.UUID `json:"id"`
	EventID        uuid.UUID `json:"event_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Status         string    `json:"status"`
	StatusCode     int       `json:"status_code"`
	DurationMs     int64     `json:"duration_ms"`
	AttemptNumber  int       `json:"attempt_number"`
	CreatedAt      time.Time `json:"created_at"`
}

// Client represents a single WebSocket connection.
type Client struct {
	conn   *websocket.Conn
	appID  string // subscribed application ID; empty = all
	send   chan []byte
	hub    *Hub
	mu     sync.Mutex
}

func (c *Client) Send(msg []byte) {
	select {
	case c.send <- msg:
	default:
		// client too slow, drop message
	}
}

// Hub manages all WebSocket clients and broadcasts.
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	logger     *zap.Logger
	mu         sync.RWMutex
}

func NewHub(logger *zap.Logger) *Hub {
	h := &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		broadcast:  make(chan []byte, 256),
		logger:     logger,
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Debug("WebSocket client connected",
				zap.String("app_id", client.appID),
				zap.Int("total", len(h.clients)),
			)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			h.logger.Debug("WebSocket client disconnected",
				zap.String("app_id", client.appID),
				zap.Int("total", len(h.clients)),
			)

		case msg := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				client.Send(msg)
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastDelivery sends a delivery update to all connected clients.
// Clients subscribed to a specific appID only receive updates for that app.
func (h *Hub) BroadcastDelivery(appID string, update *DeliveryUpdate) {
	data, err := json.Marshal(update)
	if err != nil {
		h.logger.Error("Failed to marshal delivery update", zap.Error(err))
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.appID == "" || client.appID == appID {
			client.Send(data)
		}
	}
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
