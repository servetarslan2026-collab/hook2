package ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Upgrade checks the WebSocket upgrade request and returns the upgrade handler.
func Upgrade() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only allow WebSocket upgrade
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}
		return c.Next()
	}
}

// Handle is the WebSocket connection handler.
func Handle(hub *Hub) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		logger := hub.logger
		// Read app_id from query params for scoped subscriptions
		appID := conn.Query("app_id", "")

		client := &Client{
			conn:  conn,
			appID: appID,
			send:  make(chan []byte, 64),
			hub:   hub,
		}

		hub.register <- client

		// Ensure cleanup on exit
		defer func() {
			hub.unregister <- client
			conn.Close()
		}()

		// Write pump: send messages from channel to WebSocket
		go func() {
			for msg := range client.send {
				client.mu.Lock()
				err := conn.WriteMessage(websocket.TextMessage, msg)
				client.mu.Unlock()
				if err != nil {
					break
				}
			}
		}()

		// Read pump: keep connection alive, handle pings/pongs
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					logger.Warn("WebSocket read error", zap.Error(err))
				}
				break
			}

			// Handle client messages (e.g., subscribe to specific app)
			var cmd struct {
				Action string `json:"action"`
				AppID  string `json:"app_id"`
			}
			if err := json.Unmarshal(msg, &cmd); err == nil {
				switch cmd.Action {
				case "subscribe":
					client.appID = cmd.AppID
					logger.Debug("Client subscribed to app", zap.String("app_id", cmd.AppID))
				case "unsubscribe":
					client.appID = ""
					logger.Debug("Client unsubscribed, receiving all updates")
				}
			}
		}
	})
}

// StatusHandler returns WebSocket hub status (for health/debug).
func StatusHandler(hub *Hub) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ws_clients": hub.ClientCount(),
		})
	}
}
