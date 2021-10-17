package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	models "github.com/johnythelittle/goupdateyourself/models/message"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var messages = mongoutil.DB("messages")

func checkParticipants(id primitive.ObjectID, conversationID primitive.ObjectID) (error, []models.Message) {
	var dialogue models.Dialogue
	messages.FindOne(context.TODO(), bson.M{"_id": conversationID}).Decode(&dialogue)
	var b bool = false
	for _, e := range dialogue.Participants {
		fmt.Println(e)
		if id == e {
			b = true
		}
	}
	if !b {
		return fmt.Errorf("%v", "access denied"), nil
	}
	return nil, dialogue.Messages
}
func SendMessage(c *gin.Context) {
	userId, _ := c.Get("id")
	userIdFormatted, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))
	var msg models.Message
	c.ShouldBindJSON(&msg)
	msg.Sender = userIdFormatted
	msg.Date = primitive.NewDateTimeFromTime(time.Now())
	err, msgs := checkParticipants(userIdFormatted, msg.Addressee)
	if err != nil {
		c.AbortWithStatusJSON(405, gin.H{"message": "you're not part of this group"})
	}
	msgs = append(msgs, msg)
	messages.UpdateOne(context.TODO(), bson.M{"_id": msg.Addressee}, bson.M{"$push": bson.M{"messages": msg}})
	c.JSON(200, gin.H{"messages": msgs})
}
func AddDialogue(c *gin.Context) {
	userId, _ := c.Get("id")
	userIdFormatted, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))
	var dialogue models.Dialogue
	c.ShouldBindJSON(&dialogue)
	dialogue.CreatedBy = userIdFormatted
	dialogue.Participants = append(dialogue.Participants, userIdFormatted)
	messages.InsertOne(context.TODO(), dialogue)
	c.JSON(200, gin.H{"data": dialogue})
}
func GetMyDialogues(c *gin.Context) {
	userId, _ := c.Get("id")
	userIdFormatted, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))
	var dialogues []bson.M
	cursor, err := messages.Find(context.TODO(), bson.M{"participants": bson.M{"$in": []primitive.ObjectID{userIdFormatted}}})
	if err != nil {
		c.AbortWithStatusJSON(501, gin.H{"message": "internal error on server occured"})
	}
	if err = cursor.All(context.TODO(), &dialogues); err != nil {
		log.Println("error")
	}
	fmt.Println(dialogues)

	c.JSON(200, gin.H{"dialogues": dialogues})

}
