package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	models "github.com/johnythelittle/goupdateyourself/models/profile"
	models_user "github.com/johnythelittle/goupdateyourself/models/user"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var profiles = mongoutil.DB("profile")
var user_ = mongoutil.DB("user")

func GetPublicUsers(c *gin.Context) {
	var userIds []primitive.ObjectID

	var foundUsers []models_user.User
	var foundProfiles []models.Profile

	cur, err := profiles.Find(context.TODO(), bson.D{{"is_private", false}})
	if err != nil {
		fmt.Println(err)
	}
	if err = cur.All(context.TODO(), &foundProfiles); err != nil {
		fmt.Println(err)
	}
	for _, e := range foundProfiles {
		usrId, _ := primitive.ObjectIDFromHex(e.User)
		userIds = append(userIds, usrId)
	}

	for _, usr := range userIds {
		func() {
			var result models_user.User
			user_.FindOne(context.TODO(), bson.D{{"_id", usr}}).Decode(&result)
			foundUsers = append(foundUsers, result)
		}()
	}

	fmt.Println("USER IDS", userIds)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	fmt.Println(foundUsers)

	type CombinedModel struct {
		ID      string         `json:"id"`
		URL     string         `json:"url"`
		Email   string         `json:"email"`
		Name    string         `json:"name"`
		Profile models.Profile `json:"profile"`
	}
	var arrayToReturn []CombinedModel

	for _, fu := range foundUsers {
		for _, fp := range foundProfiles {
			var combined CombinedModel
			if fp.User == fu.ID {
				combined.ID = fu.ID
				combined.Name = fu.Name
				combined.URL = fu.URLName
				combined.Email = fu.Email
				combined.Profile = fp
				arrayToReturn = append(arrayToReturn, combined)
			}

		}

	}
	c.JSON(200, gin.H{"profiles": arrayToReturn})
}
