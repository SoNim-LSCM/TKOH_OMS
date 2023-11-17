package websocket

// reference from https://github.com/gofiber/contrib/tree/main/websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/models/loginAuth"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const SubscribeTokenResponse string = `{
	"messageCode": "CONNECTION_REGISTERED",
	"userId": 1,
	"username": "user1",
	"userType" : "STAFF"
 }`

var wsConnPair = make(map[string]*websocket.Conn)

func SetupWebsocket() {
	app := fiber.New()
	// defer app.Shutdown()

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/oms/", websocket.New(func(c *websocket.Conn) {

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			// mt       int
			err      error
			loggedIn bool = false
			wsLive   bool = true
		)

		type SubscribeTokenDTO struct {
			Username  string `json:"username"`
			AuthToken string `json:"authToken"`
		}
		request := new(SubscribeTokenDTO)
		c.SetCloseHandler(func(code int, text string) error {
			delete(wsConnPair, c.RemoteAddr().String())
			wsLive = false
			return c.Close()
		})
		wsConnPair[c.RemoteAddr().String()] = c

		for wsLive {
			err = c.ReadJSON(request)
			if err == nil {
				fmt.Println(SendBoardcastMessage("123123"))
				if !loggedIn {
					if request.AuthToken == "" || request.Username == "" {

					} else {
						log.Printf("Login Username: %s , AuthToken: %s\n", request.Username, request.AuthToken)
						var response loginAuth.SubscribeTokenResponse
						err := json.Unmarshal([]byte(SubscribeTokenResponse), &response)
						errorHandler.CheckError(err, "translate string to json in wsHandler")
						err = SendBoardcastMessage(response)
						errorHandler.CheckError(err, "Error in translating message to websocket message")
						loggedIn = true
					}
				} else {
				}
			}
		}

	}))
	port := os.Getenv("WS_PORT")
	log.Println(app.Listen(":" + port))
}

func SendBoardcastMessage(msg interface{}) error {
	for addr, wsConn := range wsConnPair {
		err := wsConn.WriteJSON(msg)
		if err != nil {
			delete(wsConnPair, addr)
		}
	}
	return nil
}

func SenDirectMessage(msg interface{}) error {
	for _, wsConn := range wsConnPair {
		err := wsConn.WriteJSON(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
