package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db_tree = mongoutil.DB("trees")

type Tree struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	User      string             `json:"user" bson:"user"`
	Branches  []Branch           `json:"branches" bson:"branches"`
	IsPrivate bool               `json:"isPrivate" bson:"isPrivate"`
}
type Branch struct {
	Nodes     []Node `json:"nodes" bson:"nodes"`
	IsPrivate bool   `json:"isPrivate" bson:"isPrivate"`
	Name      string `json:"nameOfBranch" bson:"nameOfBranch"`
}
type BranchFromCli struct {
	OfTree    string `json:"belongsToTree"`
	IsPrivate bool   `json:"isPrivate" bson:"isPrivate"`
	Name      string `json:"nameOfBranch" bson:"nameOfBranch"`
}
type Node struct {
	Value string `bson:"valueOfNode"`
}
type NodeInfoFromCli struct {
	NumberOfBranch int    `json:"numberOfBranch"`
	IdOfTree       string `json:"idOfTree"`
	Value          string `json:"valueOfNode"`
}

func AddTree(c *gin.Context) {
	userId, _ := c.Get("id")
	var tree Tree
	c.ShouldBindJSON(&tree)
	result, _ := db_tree.InsertOne(context.TODO(), bson.D{{"user", userId}, {"isPrivate", tree.IsPrivate}, {"name", tree.Name}, {"branches", []Branch{}}})
	c.JSON(200, gin.H{
		"result": result.InsertedID,
	})
}

func AddBranch(c *gin.Context) {

	var branchInfo BranchFromCli
	c.ShouldBindJSON(&branchInfo)
	addBranch := Branch{Name: branchInfo.Name, Nodes: []Node{}, IsPrivate: branchInfo.IsPrivate}
	idOfTree, _ := primitive.ObjectIDFromHex(branchInfo.OfTree)
	db_tree.FindOneAndUpdate(context.TODO(), bson.M{"_id": idOfTree}, bson.M{"$push": bson.M{"branches": addBranch}})
	var res Tree
	db_tree.FindOne(context.TODO(), bson.M{"_id": idOfTree}).Decode(&res)
	fmt.Println(res)
}

func AddNode(c *gin.Context) {
	var nodeInfo NodeInfoFromCli
	c.ShouldBindJSON(&nodeInfo)

	fmt.Println(nodeInfo.NumberOfBranch, nodeInfo.IdOfTree, nodeInfo.Value)
	idOfTree, _ := primitive.ObjectIDFromHex(nodeInfo.IdOfTree)

	filter := bson.M{"_id": idOfTree}
	toString := strconv.Itoa(nodeInfo.NumberOfBranch)
	searchString := "branches." + toString + ".nodes"

	update := bson.M{"$push": bson.M{searchString: bson.M{"valueOfNode": nodeInfo.Value}}}

	db_tree.UpdateOne(context.TODO(), filter, update)

}

func GetTree(c *gin.Context) {
	userId, _ := c.Get("id")
	var trees []Tree
	cur, _ := db_tree.Find(context.TODO(), bson.M{"user": userId})
	cur.All(context.TODO(), &trees)
	c.JSON(200, gin.H{"result": trees})
}
