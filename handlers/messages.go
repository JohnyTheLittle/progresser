package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	models "github.com/johnythelittle/goupdateyourself/models/message"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var messages = mongoutil.DB("messages")

func GetMyMessages(c *gin.Context) {
	userId, _ := c.Get("id")
	formattedUserID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))
	var msgs []models.Message
	var msgs_s []models.Message
	var msgs_r []models.Message

	list_sent, err := messages.Find(context.TODO(), bson.D{{"sender", formattedUserID}})
	if err != nil {
		fmt.Println(err)
	}
	if err = list_sent.All(context.TODO(), &msgs_s); err != nil {
		fmt.Println(err)
	}

	list_recieved, err := messages.Find(context.TODO(), bson.D{{"addressee", formattedUserID}})
	if err != nil {
		fmt.Println(err)
	}
	if err = list_recieved.All(context.TODO(), &msgs_r); err != nil {
		fmt.Println(err)
	}

	msgs = append(msgs, msgs_s...)
	msgs = append(msgs, msgs_r...)

	bs, _ := json.Marshal(msgs)
	fmt.Println(string(bs))
	c.JSON(200, gin.H{"data": msgs})
}

func SendMessage(c *gin.Context) {
	userId, _ := c.Get("id")

	formattedUserID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))

	var msg models.Message

	c.ShouldBindJSON(&msg)
	msg.Sender = formattedUserID
	msg.Date = primitive.NewDateTimeFromTime(time.Now())

	bs, _ := json.Marshal(msg)
	fmt.Println(string(bs))

	messages.InsertOne(context.TODO(), msg)

	c.JSON(200, gin.H{"data": msg})

}
