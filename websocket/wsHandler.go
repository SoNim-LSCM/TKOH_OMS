package websocket

// reference from https://github.com/gofiber/contrib/tree/main/websocket

import (
	"encoding/json"
	"log"
	"os"
	"tkoh_oms/errors"
	"tkoh_oms/models/loginAuth"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const SubscribeTokenResponse string = `{
	"messageCode": "CONNECTION_REGISTERED",
	"userId": 1,
	"username": "user1",
	"userType" : "STAFF"
 }`

var wsObject *websocket.Conn

func SetupWebsocket() {
	app := fiber.New()
	defer app.Shutdown()

	app.Get("/oms/", websocket.New(func(c *websocket.Conn) {

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt       int
			err      error
			loggedIn bool = false
		)

		type SubscribeTokenDTO struct {
			Username  string `json:"username"`
			AuthToken string `json:"authToken"`
		}
		wsObject = c
		request := new(SubscribeTokenDTO)

		for {
			if !loggedIn {
				if err = wsObject.ReadJSON(request); err != nil || (request.AuthToken == "" || request.Username == "") {
					ret := []byte(err.Error())
					log.Println("request:", err)
					if err = wsObject.WriteMessage(mt, ret); err != nil {
						log.Println("write:", err)
					}
				} else {
					log.Printf("Login Username: %s , AuthToken: %s\n", request.Username, request.AuthToken)
					var response loginAuth.SubscribeTokenResponse
					err := json.Unmarshal([]byte(SubscribeTokenResponse), &response)
					errors.CheckError(err, "translate string to json in wsHandler")
					SendMessage(response)
					loggedIn = true
				}
			}
		}

	}))
	port := os.Getenv("WS_PORT")
	log.Println(app.Listen(":" + port))
}

func SendMessage(msg interface{}) error {
	err := wsObject.WriteJSON(msg)
	if err != nil {
		log.Println("read:", err)
	}
	return err
}
