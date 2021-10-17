package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	models "github.com/johnythelittle/goupdateyourself/models/profile"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
)

var profile = mongoutil.DB("profile")

func GetMyProfile(c *gin.Context) {
	userId, _ := c.Get("id")
	var myProfile models.Profile
	profile.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&myProfile)
	c.JSON(200, gin.H{"data": myProfile})
}

func GetProfile(c *gin.Context) {
	//ADD BLACK LILST
	//blacklist := mongoutil.DB("bl")
	requiredId := c.Query("id")
	var profile_ models.Profile

	profile.FindOne(context.TODO(), bson.M{"user_id": requiredId}).Decode(&profile_)

	//blacklist.FindOne(context.TODO(), bson.D{{}})

	if profile_.IsPrivate {
		c.AbortWithStatusJSON(405, gin.H{"message": "ACCESS DENIED"})
	} else {
		c.JSON(200, gin.H{"data": profile_})
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

	profile.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&profile_)

	if profile_.User == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.D{{"$set", bson.D{{"age", data.Age}}}})
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

	profile.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&profile_)

	listOfEdus := profile_.Education

	listOfEdus = append(listOfEdus, data.Education)

	fmt.Println(listOfEdus)
	if profile_.User == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.D{{"$set", bson.D{{"education", listOfEdus}}}})
		c.JSON(200, gin.H{"message": "success"})
	} else {
		c.AbortWithStatusJSON(405, gin.H{"message": "ACCESS DENIED"})
	}
}
func DeleteEducation(c *gin.Context) {
	userId, _ := c.Get("id")
	type DeletedEducation struct {
		DeletedEducation int `json:"number"`
	}
	var data DeletedEducation
	c.ShouldBindJSON(&data)
	var profile_ models.Profile
	profile.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&profile_)
	listOfEdus := profile_.Education
	updatedlistOfEdus := append(listOfEdus[:data.DeletedEducation], listOfEdus[data.DeletedEducation+1:]...)
	if profile_.User == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.D{{"$set", bson.D{{"education", updatedlistOfEdus}}}})
		c.JSON(200, gin.H{"data": "success"})
	} else {
		c.AbortWithStatusJSON(405, gin.H{"message": "ACCESS DENIED"})
	}
}
func DeletePerk(c *gin.Context) {
	userId, _ := c.Get("id")
	type DeletedPerk struct {
		DeletedPerk int `json:"number"`
	}
	var data DeletedPerk
	c.ShouldBindJSON(&data)
	var profile_ models.Profile
	profile.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&profile_)
	listOfPerks := profile_.Perks
	updatedListOfPerks := append(listOfPerks[:data.DeletedPerk], listOfPerks[data.DeletedPerk+1:]...)
	if profile_.User == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.D{{"$set", bson.D{{"perks", updatedListOfPerks}}}})
		c.JSON(200, gin.H{"data": "success"})
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

	profile.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&profile_)

	listOfPerks := profile_.Perks

	listOfPerks = append(listOfPerks, data.Perk)

	if profile_.User == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.D{{"$set", bson.D{{"perks", listOfPerks}}}})
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

	profile.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&profile_)
	if profile_.User == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.D{{"$set", bson.D{{"description", data}}}})
		c.JSON(200, gin.H{"message": "success"})
	} else {
		c.AbortWithStatus(405)
	}
}

func TogglePrivacy(c *gin.Context) {
	userId, _ := c.Get("id")
	var profile_ models.Profile

	profile.FindOne(context.TODO(), bson.M{"user_id": userId}).Decode(&profile_)

	if profile_.User == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.D{{"$set", bson.D{{"is_private", !profile_.IsPrivate}}}})
		c.JSON(200, gin.H{"message": !profile_.IsPrivate})
	} else {
		c.AbortWithStatusJSON(405, gin.H{"message": "ACCESS DENIED"})
	}
}

func SetProfileData(c *gin.Context) {
	//WHEN REGISTER APPEND STRING NOT AN ID!!!!!!!!!!!!!

	userId, _ := c.Get("id")
	var profileInfo models.Profile
	c.ShouldBindJSON(&profileInfo)
	fmt.Println(profileInfo)
	var requiredProfile models.Profile

	profile.FindOne(context.TODO(), bson.D{{"user_id", userId}}).Decode(&requiredProfile)

	fmt.Println("BELONGS:::", requiredProfile)

	if requiredProfile.User == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.D{{"$set", bson.D{{"age", profileInfo.Age}, {"education", profileInfo.Education}, {"perks", profileInfo.Perks}, {"description", profileInfo.SelfRepresentation}, {"pronounce", profileInfo.Pronounce}, {"is_private", profileInfo.IsPrivate}}}})
		c.JSON(200, gin.H{"message": "success"})
	} else {
		c.AbortWithStatus(405)
	}

}
func SetPronounce(c *gin.Context) {
	userId, _ := c.Get("id")
	type Pronounce struct {
		Value string `json:"pronounce"`
	}
	var pronounce Pronounce
	var requiredProfile models.Profile
	c.ShouldBindJSON(&pronounce)
	profile.FindOne(context.TODO(), bson.D{{"user_id", userId}}).Decode(&requiredProfile)

	if requiredProfile.User == userId {
		profile.UpdateOne(context.TODO(), bson.M{"user_id": userId}, bson.D{{"$set", bson.D{{"pronounce", pronounce.Value}}}})
		c.JSON(200, gin.H{"message": "success"})
	} else {
		c.AbortWithStatus(405)
	}
}
