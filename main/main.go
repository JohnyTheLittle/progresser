package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/johnythelittle/goupdateyourself/configs"
	"github.com/johnythelittle/goupdateyourself/handlers"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
)

func main() {
	c, err := configs.LoadConfig("../")
	if err != nil {
		log.Panic("BAD CONFIG")
	}
	gin.ForceConsoleColor()
	r := gin.Default()

	mongoutil.DB("")

	v1 := r.Group("/public")
	{
		v1.POST("/login", handlers.Login)
		v1.POST("/register", handlers.CreateUser)
	}

	v2 := r.Group("/private", handlers.CheckUser)
	{
		go v2.POST("/addTree", handlers.AddTree)
		go v2.POST("/addBranch", handlers.AddBranch)
		go v2.POST("/addNode", handlers.AddNode)
		go v2.GET("/myTrees", handlers.GetTree)

		go v2.POST("/personalValues", handlers.AddPersonalvalue)
	}

	r.Run(c.Port)
}
