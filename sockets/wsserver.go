package wsserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/johnythelittle/goupdateyourself/configs"
	models "github.com/johnythelittle/goupdateyourself/models/message"
	modelsUser "github.com/johnythelittle/goupdateyourself/models/user"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var user = mongoutil.DB("user")
var dialogue = mongoutil.DB("dialogues")
var config, _ = configs.LoadConfig("../")

type ClientsManager struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

var Manager = ClientsManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Clients:    make(map[string]*Client),
	Unregister: make(chan *Client),
}

func (manager *ClientsManager) Start() {
	for {
		fmt.Println("INITIALIZING WEBSOCKET SERVER")
		select {
		case conn := <-Manager.Register:
			Manager.Clients[conn.ID] = conn
			jsonMessage, _ := json.Marshal(&models.Message{Text: "Conversation has been started"})
			conn.Send <- jsonMessage

		case conn := <-Manager.Unregister:
			if _, ok := Manager.Clients[conn.ID]; ok {
				jsonMessage, _ := json.Marshal(&models.Message{Text: "Conversation has been stopped"})
				conn.Send <- jsonMessage
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case message := <-Manager.Broadcast:
			MessageStruct := models.Message{}
			json.Unmarshal(message, &MessageStruct)
			for id, conn := range Manager.Clients {
				fmt.Println("here is the problem:")
				fmt.Printf("id: %v\n", id)
				fmt.Printf("conn: %v\n", conn)
				if id != conn.ID {
					continue
				}
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}

		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println("message sent to a pool", c.ID)
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		log.Printf("message read to client: %s", string(message))
		Manager.Broadcast <- message
	}
}

func ChatHandler(c *gin.Context) {
	var stringToken = c.Query("token")
	var toWhom = c.Query("adressee")
	var userInfo modelsUser.User
	var dialogue_ models.Dialogue
	fmt.Println(toWhom, stringToken)
	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("ERROR ERROR")
		}
		return []byte(config.Secret), nil
	})

	if err != nil {
		fmt.Println(err)
		c.AbortWithError(405, err)
	}

	if token.Valid {
		mapstructure.Decode(token.Claims, &userInfo)
		err := user.FindOne(context.TODO(), bson.D{{"email", userInfo.Email}}).Decode(&userInfo)
		if err != nil {
			fmt.Println("ERROR CANNOT FIND USER", err)
			c.AbortWithError(400, err)
		}
		uid := userInfo.ID
		conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			http.NotFound(c.Writer, c.Request)
			return
		}

		client := &Client{
			ID:     dialogue_.ID,
			Socket: conn,
			Send:   make(chan []byte),
		}
		Manager.Register <- client
		_, textOfMessage, _ := client.Socket.ReadMessage()
		var message models.Message
		message.Text = string(textOfMessage)
		message.From = uid
		message.To = toWhom
		message.Date = primitive.DateTime(time.Now().Unix())
		err = dialogue.FindOne(context.TODO(), bson.D{{
			"participants", bson.D{{"$all", bson.A{uid, toWhom}}}}}).Decode(&dialogue_)
		if dialogue_.CreatedBy.IsZero() {
			dialogue.InsertOne(context.TODO(), bson.D{{"participants", []string{userInfo.ID, toWhom}}, {"creator", uid}, {"messages", bson.A{message}}})
		} else {
			messages := dialogue_.Messages
			messages = append(messages, message)
			dialogue.UpdateOne(context.TODO(), bson.D{{
				"participants", bson.D{{"$all", bson.A{uid, toWhom}}}}}, bson.M{"$set": bson.M{"messages": messages}})
		}
		go client.Read()
		go client.Write()

	}
}
