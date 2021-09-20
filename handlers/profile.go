package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	models "github.com/johnythelittle/goupdateyourself/models/profile"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var profile = mongoutil.DB("profile")

func GetProfile(c *gin.Context) {
	userId, _ := c.Get("id")
	requiredId := c.Query("id")
	var profile_ models.Profile
	userIDFormated, _ := primitive.ObjectIDFromHex(requiredId)
	profile.FindOne(context.TODO(), bson.M{"user": userIDFormated}).Decode(&profile_)
	if userId == requiredId {
		c.JSON(200, gin.H{"data": profile_})
	} else {
		if profile_.IsPrivate {
			c.AbortWithStatusJSON(405, gin.H{"message": "ACCESS DENIED"})
		} else {
			c.JSON(200, gin.H{"data": profile_})
		}
	}

}

func SetAge(c *gin.Context) {
	userId, _ := c.Get("id")

	type Age struct {
		Age int `json:"age"`
	}
	var data Age

	c.ShouldBindJSON(&data)
	var profile_ models.Profile

	userIDFormated, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))

	profile.FindOne(context.TODO(), bson.M{"user": userIDFormated}).Decode(&profile_)

	if profile_.Belongs == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user": userIDFormated}, bson.D{{"$set", bson.D{{"age", data.Age}}}})
		c.JSON(200, gin.H{"message": "success"})
	} else {
		c.AbortWithStatusJSON(405, gin.H{"message": "ACCESS DENIED"})
	}

}

func AddEducation(c *gin.Context) {
	userId, _ := c.Get("id")
	type Education struct {
		Education models.Education `json:"edu"`
	}
	var data Education
	c.ShouldBindJSON(&data)
	var profile_ models.Profile
	userIDFormated, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))

	profile.FindOne(context.TODO(), bson.M{"user": userIDFormated}).Decode(&profile_)

	listOfEdus := profile_.Education

	listOfEdus = append(listOfEdus, data.Education)

	fmt.Println(listOfEdus)
	if profile_.Belongs == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user": userIDFormated}, bson.D{{"$set", bson.D{{"education", listOfEdus}}}})
		c.JSON(200, gin.H{"message": "success"})
	} else {
		c.AbortWithStatusJSON(405, gin.H{"message": "ACCESS DENIED"})
	}
}

func AddPerk(c *gin.Context) {
	userId, _ := c.Get("id")
	type Perk struct {
		Perk models.Perk `json:"perk"`
	}
	var data Perk
	c.ShouldBindJSON(&data)
	var profile_ models.Profile
	userIDFormated, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))

	profile.FindOne(context.TODO(), bson.M{"user": userIDFormated}).Decode(&profile_)

	listOfPerks := profile_.Perks

	listOfPerks = append(listOfPerks, data.Perk)

	if profile_.Belongs == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user": userIDFormated}, bson.D{{"$set", bson.D{{"perks", listOfPerks}}}})
		c.JSON(200, gin.H{"message": "success"})
	} else {
		c.AbortWithStatusJSON(405, gin.H{"message": "ACCESS DENIED"})
	}
}

func AddSelfRepresentation(c *gin.Context) {
	userId, _ := c.Get("id")

	var data models.SelfRepresentation

	c.ShouldBindJSON(&data)

	var profile_ models.Profile

	userIDFormated, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))

	profile.FindOne(context.TODO(), bson.M{"user": userIDFormated}).Decode(&profile_)
	if profile_.Belongs == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user": userIDFormated}, bson.D{{"$set", bson.D{{"description", data.Text}}}})
		c.JSON(200, gin.H{"message": "success"})
	} else {
		c.AbortWithStatus(405)
	}
}

func TogglePrivacy(c *gin.Context) {
	userId, _ := c.Get("id")
	userIDFormated, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))
	var profile_ models.Profile

	profile.FindOne(context.TODO(), bson.M{"user": userIDFormated}).Decode(&profile_)

	if profile_.Belongs == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user": userIDFormated}, bson.D{{"$set", bson.D{{"is_private", !profile_.IsPrivate}}}})
		c.JSON(200, gin.H{"message": !profile_.IsPrivate})
	} else {
		c.AbortWithStatusJSON(405, gin.H{"message": "ACCESS DENIED"})
	}
}

func SetProfileData(c *gin.Context) {
	userId, _ := c.Get("id")

	var profileInfo models.Profile

	c.ShouldBindJSON(&profileInfo)

	userIDFormated, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))

	var requiredProfile models.Profile

	profile.FindOne(context.TODO(), bson.M{"user": userIDFormated}).Decode(&requiredProfile)
	if requiredProfile.Belongs == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user": userIDFormated}, bson.D{{"$set", bson.D{{"age", profileInfo.Age}, {"education", profileInfo.Education}, {"perks", profileInfo.Perks}, {"description", profileInfo.SelfRepresentation.Text}, {"pronounce", profileInfo.Pronounce}, {"is_private", profileInfo.IsPrivate}}}})
		c.JSON(200, gin.H{"message": "success"})
	} else {
		c.AbortWithStatus(405)
	}

}
