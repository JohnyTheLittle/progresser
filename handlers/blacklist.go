package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	models "github.com/johnythelittle/goupdateyourself/models/blacklist"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var bl = mongoutil.DB("blacklists")

func AddToBlackList(c *gin.Context) {
	userId, _ := c.Get("id")
	userIdFormatted, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))
	var existingList models.Blacklist
	type PersonToAdd struct {
		ID primitive.ObjectID `json:"id"`
	}
	var personAdded PersonToAdd
	c.ShouldBindJSON(&personAdded)
	bl.FindOne(context.TODO(), bson.D{{"user", userIdFormatted}}).Decode(&existingList)
	if existingList.Belongs.IsZero() {
		fmt.Println("it is")
		bl.InsertOne(context.TODO(), bson.D{{"user", userIdFormatted}, {"blocked_persons", []primitive.ObjectID{personAdded.ID}}})
	} else {
		existingList.ListOfPersons = append(existingList.ListOfPersons, personAdded.ID)
		bl.UpdateOne(context.TODO(), bson.D{{"user", userIdFormatted}}, bson.D{{"$set", bson.D{{"blocked_persons", existingList.ListOfPersons}}}})
		fmt.Println(existingList.ListOfPersons)
	}
	c.JSON(200, gin.H{"message": "success"})
}

func RemoveFromBlackList(c *gin.Context) {
	userId, _ := c.Get("id")
	userIdFormatted, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))
	var existingList models.Blacklist
	type PersonToRemove struct {
		ID primitive.ObjectID `json:"id"`
	}
	var removedPerson PersonToRemove
	c.ShouldBindJSON(&removedPerson)
	bl.FindOne(context.TODO(), bson.D{{"user", userIdFormatted}}).Decode(&existingList)

	var index int
	for i, e := range existingList.ListOfPersons {
		if e == removedPerson.ID {
			index = i
			break
		}
	}
	updatedList := append(existingList.ListOfPersons[:index], existingList.ListOfPersons[index+1:]...)
	bl.UpdateOne(context.TODO(), bson.D{{"user", userIdFormatted}}, bson.D{{"$set", bson.D{{"blocked_persons", updatedList}}}})
}
