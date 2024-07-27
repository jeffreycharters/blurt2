package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"strconv"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jeffreycharters/blurt/db"
	"github.com/jeffreycharters/blurt/ent"
	"github.com/jeffreycharters/blurt/ent/blurt"
	"github.com/jeffreycharters/blurt/ent/user"
)

// Add more data to this type if needed
type client struct {
	isClosing bool
	mu        sync.Mutex
}

var clients = make(map[*websocket.Conn]*client)
var register = make(chan *websocket.Conn)
var broadcast = make(chan []byte)
var unregister = make(chan *websocket.Conn)

func blurtHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = &client{}
			log.Println("connection registered")

		case message := <-broadcast:
			log.Println("message broadcasting:", string(message))
			// Send the message to all clients
			for connection, c := range clients {
				go func(connection *websocket.Conn, c *client) { // send to each client in parallel so we don't block on a slow client
					c.mu.Lock()
					defer c.mu.Unlock()
					if c.isClosing {
						return
					}
					if err := connection.WriteMessage(websocket.TextMessage, message); err != nil {
						c.isClosing = true
						log.Println("write error:", err)

						connection.WriteMessage(websocket.CloseMessage, []byte{})
						connection.Close()
						unregister <- connection
					}
				}(connection, c)
			}

		case connection := <-unregister:
			// Remove the client from the hub
			delete(clients, connection)

			log.Println("connection unregistered")
		}
	}
}

func main() {
	app := fiber.New()

	app.Use(cors.New())

	client := db.Client()
	defer client.Close()

	app.Post("/api/v1/users/:username", func(c *fiber.Ctx) error {
		username, err := url.QueryUnescape(c.Params("username"))
		if err != nil {
			return c.SendStatus(500)
		}

		user, err := client.User.Query().Where(user.Username(username)).First(context.Background())
		if ent.IsNotFound(err) {
			user, err = client.User.Create().SetUsername(username).Save(context.Background())
			if err != nil {
				return c.SendStatus(500)
			}
		} else if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(user)
	})

	app.Get("/api/v1/blurts", func(c *fiber.Ctx) error {
		offset, err := strconv.Atoi(c.Query("offset", "0"))
		if err != nil {
			offset = 0
		}

		loadCount, err := strconv.Atoi(c.Query("count", "25"))
		if err != nil {
			loadCount = 25
		}

		log.Println("offset:", offset, "loadCount:", loadCount)

		blurts, err := client.Blurt.Query().
			WithAuthor().
			WithLiks().
			Order(ent.Desc((blurt.FieldCreateTime))).
			Offset(offset).
			Limit(loadCount).
			All(context.Background())
		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(blurts)
	})

	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	go blurtHub()

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		defer func() {
			unregister <- c
			c.Close()
		}()

		// Register the client
		register <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// Broadcast the received message
				ret, err := db.ParseMessage(message)
				if err != nil {
					c.WriteJSON(db.WrapError(err))
					return
				}
				broadcast <- ret
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))

	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))
}
