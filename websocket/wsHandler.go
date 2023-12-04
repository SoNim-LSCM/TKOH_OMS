package websocket

// reference from https://github.com/gofiber/contrib/tree/main/websocket

import (
	"log"
	"os"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/SoNim-LSCM/TKOH_OMS/models/loginAuth"
	"github.com/SoNim-LSCM/TKOH_OMS/utils"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var wsConnPair = make(map[int]*websocket.Conn)

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
		// c.SetCloseHandler(func(code int, text string) error {
		// 	delete(wsConnPair, c.RemoteAddr().String())
		// 	// delete(wsConnPair, c.RemoteAddr().String())
		// 	wsLive = false
		// 	return c.Close()
		// })
		// wsConnPair[c.RemoteAddr().String()] = c

		for wsLive {
			err = c.ReadJSON(request)
			if err == nil {
				var response interface{}
				if !loggedIn {
					response = loginAuth.FailSubscribeTokenResponse{MessageCode: "FAILED", FailReason: "User Not Found"}
					if request.AuthToken != "" || request.Username != "" {
						log.Printf("Login Username: %s , AuthToken: %s\n", request.Username, request.AuthToken)
						isValid, err := utils.ValidateJwtToken(request.AuthToken)
						if isValid && err == nil {
							claims, err := utils.GetDetailsJwtToken(request.AuthToken)
							if claims.Username == request.Username && err == nil {
								response = loginAuth.SubscribeTokenResponse{MessageCode: "CONNECTION_REGISTERED", UserID: claims.UserId, Username: claims.Username, UserType: claims.UserType}
								c.SetCloseHandler(func(code int, text string) error {
									delete(wsConnPair, claims.UserId)
									wsLive = false
									return c.Close()
								})
								wsConnPair[claims.UserId] = c
								loggedIn = true
							} else {
								response = loginAuth.FailSubscribeTokenResponse{MessageCode: "FAILED", FailReason: "Broken token / token unmatch with user"}
							}
						} else {
							response = loginAuth.FailSubscribeTokenResponse{MessageCode: "FAILED", FailReason: "Invalid token"}
						}
					} else {
						response = loginAuth.FailSubscribeTokenResponse{MessageCode: "FAILED", FailReason: "Invalid Input"}
					}
					err = c.WriteJSON(response)
					errorHandler.CheckError(err, "Error in write json to websocket")
				} else {
					c.WriteJSON(request)
				}
			}
		}

	}))
	port := os.Getenv("WS_PORT")
	log.Println(app.Listen(":" + port))
}

func SendBoardcastMessage(msg interface{}) error {
	for _, wsConn := range wsConnPair {
		err := wsConn.WriteJSON(msg)
		errorHandler.CheckError(err, "Send Websocket Message Failed")
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
