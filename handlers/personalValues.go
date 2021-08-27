package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
)

var db_personal_values = mongoutil.DB("personalValues")

type PersonalValues struct {
	UserID    string            `bson:"user" json:"user"`
	IsPrivate bool              `bson:"isPrivate" json:"isPrivate"`
	Education []Education       `bson:"education" json:"education"`
	Perks     []IndividualPerks `bson:"perks" json:"perks"`
}
type Education struct {
	Education string `bson:"edu" json:"edu"`
	Degree    string `bson:"degree" json:"degree"`
}
type IndividualPerks struct {
	PerkName        string `bson:"perkName" json:"perkName"`
	PerkDescription string `bson:"perkDescription" json:"perkDescription"`
}

func AddPersonalvalue(c *gin.Context) {
	userId, _ := c.Get("id")
	userName, _ := c.Get("username")
	var personalValues PersonalValues
	c.ShouldBindJSON(&personalValues)
	db_personal_values.InsertOne(context.TODO(), bson.D{{"user", userId}, {"isPrivate", personalValues.IsPrivate}, {"education", personalValues.Education}, {"perks", personalValues.Perks}})
	c.JSON(200, gin.H{"uid": userId, "username": userName})
}
