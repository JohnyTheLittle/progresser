package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	models "github.com/johnythelittle/goupdateyourself/models/branch"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
)

var branch = mongoutil.DB("branches")

func CreateBranch(c *gin.Context) {
	var bpfc models.Branch
	//branch params from client side
	user_id, _ := c.Get("id")
	err := c.ShouldBindJSON(&bpfc)
	if err == nil {
		res, err := branch.InsertOne(context.TODO(), bson.D{{"user", user_id}, {"name_of_branch", bpfc.Name}, {"books", bpfc.Books}, {"projects", bpfc.Projects}, {"imrovements", bpfc.Improvement}, {"is_private", true}, {"video_courses", bpfc.VideoCourses}}, nil)
		if err != nil {
			c.JSON(500, gin.H{
				"err": err,
			})
		} else {
			c.JSON(500, gin.H{
				"data": res,
			})
		}
	}

}
func GetAllBranchesOfUser(c *gin.Context) {
	var usersBranches []models.Branch
	user_id, _ := c.Get("id")
	result, err := branch.Find(context.TODO(), bson.D{{"user", user_id}})
	if err == nil {
		result.Decode(&usersBranches)
	}

	result.All(context.TODO(), &usersBranches)
	fmt.Println(usersBranches)
	c.JSON(200, usersBranches)
}

func RenameBranch(c *gin.Context) {
	user_id, _ := c.Get("id")

	type NewNameJson struct {
		NewName string `json:"newName"`
		ID      string `json:"id"`
	}

	var new_name NewNameJson
	err := c.ShouldBindJSON(&new_name)
	id_formatted, _ := primitive.ObjectIDFromHex(new_name.ID)

	func(id primitive.ObjectID) {
		var branch_ models.Branch
		err := branch.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&branch_)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(405)
		}
		if user_id != branch_.Belongs {
			fmt.Println("USER TRIED CHANGE NAME OF BRANCH WHICH DOESENT BELONG IT")
			c.AbortWithStatus(405)
		}
	}(id_formatted)

	if err != nil {
		fmt.Println("something wrong during updating name of branch")
		c.AbortWithStatus(501)
	}
	branch.UpdateOne(context.TODO(), bson.M{"_id": id_formatted}, bson.D{{"$set", bson.D{{"name_of_branch", new_name.NewName}}}})
	c.JSON(200, gin.H{"result": "SUCCESS"})
}

func AppendNewElementToBooks(c *gin.Context) {
	user_id, _ := c.Get("id")
	type NewBook struct {
		Book models.Book `json:"appended_book"`
		ID   string      `json:"ID"`
	}
	var new_book NewBook
	c.ShouldBindJSON(&new_book)
	id_formatted, _ := primitive.ObjectIDFromHex(new_book.ID)
	var branch_ models.Branch
	func(id primitive.ObjectID) {
		err := branch.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&branch_)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(405)
		}
		if user_id != branch_.Belongs {
			fmt.Println("USER TRIED TO APPEND BOOK TO NOT ITS OWN BRANCH")
			c.AbortWithStatus(405)
		}
	}(id_formatted)
	func(book models.Book) {
		fmt.Println(branch_.Books)
		fmt.Println(book)
		updatedList := append(branch_.Books, book)
		fmt.Println(updatedList)
		branch.UpdateOne(context.TODO(), bson.M{"_id": id_formatted}, bson.D{{"$set", bson.D{{"books", updatedList}}}})
		c.JSON(200, gin.H{"result": "SUCCESS"})
	}(new_book.Book)
}
func DeleteElementFromBooks(c *gin.Context) {
	user_id, _ := c.Get("id")
	type DeletedBook struct {
		Num int    `json:"number_of_deleted_book"`
		ID  string `json:"ID"`
	}
	var deletedBook DeletedBook
	c.ShouldBindJSON(&deletedBook)
	id_formatted, _ := primitive.ObjectIDFromHex(deletedBook.ID)
	var branch_ models.Branch
	func(id primitive.ObjectID) {
		err := branch.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&branch_)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(405)
		}
		if user_id != branch_.Belongs {
			fmt.Println("USER TRIED TO DELETE BOOK FROM NOT ITS OWN BRANCH")
			c.AbortWithStatus(405)
		}
	}(id_formatted)
	func(numberOfBook int) {
		updatedList := append(branch_.Books[:numberOfBook], branch_.Books[numberOfBook+1:]...)
		branch.UpdateOne(context.TODO(), bson.M{"_id": id_formatted}, bson.D{{"$set", bson.D{{"books", updatedList}}}})
		c.JSON(200, gin.H{"result": "SUCCESS"})
	}(deletedBook.Num)
}

func AddVideoCourse(c *gin.Context) {
	user_id, _ := c.Get("id")
	type AddedVideoCourse struct {
		VideoCourse models.VideoCourse `json:"video_course"`
		ID          string             `json:"ID"`
	}
	var video_course AddedVideoCourse
	var branch_ models.Branch
	c.ShouldBindJSON(&video_course)
	id_formatted, _ := primitive.ObjectIDFromHex(video_course.ID)

	func(id primitive.ObjectID) {
		err := branch.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&branch_)
		if err != nil {
			fmt.Println("THERE IS NO SUCH BRANCH")
			c.AbortWithStatus(404)
		}
		if user_id != branch_.Belongs {
			fmt.Println("ATTEMPTED TO ADD VIDEOCOURSE TO SOMEONE ELSES BRANCH")
			c.AbortWithStatus(405)
		}
	}(id_formatted)

	func(videoCourse models.VideoCourse) {
		updatedList := append(branch_.VideoCourses, videoCourse)
		branch.UpdateOne(context.TODO(), bson.M{"_id": id_formatted}, bson.D{{"$set", bson.D{{"video_courses", updatedList}}}})
		c.JSON(200, gin.H{"result": "SUCCESS"})
	}(video_course.VideoCourse)
}
