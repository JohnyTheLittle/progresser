package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	models "github.com/johnythelittle/goupdateyourself/models/message"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
)

var dialogues = mongoutil.DB("dialogues")

type ClientManager struct {
}

func SendMessage(c *gin.Context) {
	senderID, _ := c.Get("id")
	var message_ models.Message
	var dialogue_ models.Dialogue
	senderIDFormatted := fmt.Sprintf("%v", senderID)
	c.ShouldBindJSON(&message_)
	dialogues.FindOne(context.TODO(), bson.D{{"participants", bson.M{"$in": []string{senderIDFormatted, message_.To}}}}).Decode(&dialogue_)
	fmt.Printf("dialogue_: %v\n", dialogue_)
	if dialogue_.ID == "" {
		fmt.Println("here we are. its empty")
		dialogues.InsertOne(context.TODO(), bson.D{{"participants", [2]string{senderIDFormatted, message_.To}}, {"messages", message_}})
	}

}
