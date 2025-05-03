package main

import (
	"log"
	"fmt"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/template/django/v3"
)

func main() {
	log.Println("Starting...")
	engine := django.New("./views", ".django")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
}

func setupRoutes(app *fiber.App) {
	setupWs(app)
	setupIndex(app)
}

type WebSocketMessage struct {
	ChatMessage string
}

func setupWs(app *fiber.App) {
	app.Use("/chat", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/chat", websocket.New(func(c *websocket.Conn) {
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			res := WebSocketMessage{}
			err := json.Unmarshal([]byte(msg), &res)
			if err != nil {
				log.Println("Big f")
				log.Println("%s", err)
			}

			log.Printf("msg: %s", res.ChatMessage)

			echo := fmt.Sprintf("<div id=chat-box>%s</div>", res.ChatMessage)

			if err = c.WriteMessage(mt, []byte(echo)); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))
}

func setupIndex(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"TheBestest": "Monika",
		})
	})
}
