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

	publicRoutes := r.Group("/public")
	{
		go publicRoutes.POST("/login", handlers.Login)
		go publicRoutes.POST("/register", handlers.CreateUser)
		go publicRoutes.GET("/user/id", handlers.GetUser)
		go publicRoutes.GET("/get_url", handlers.GetUrl)
		//publicRoutes.GET("/user/:user_url")
	}

	privateRoutes := r.Group("/private", handlers.CheckUser)
	{
		go privateRoutes.POST("/makeBranch", handlers.CreateBranch)
		go privateRoutes.GET("/getBranches", handlers.GetAllBranchesOfUser)
		go privateRoutes.PUT("/renameBranch", handlers.RenameBranch)
		go privateRoutes.PUT("/appendBook", handlers.AppendNewElementToBooks)
		go privateRoutes.PUT("/deleteBook", handlers.DeleteElementFromBooks)
		go privateRoutes.PUT("/appendVideoCourse", handlers.AddVideoCourse)

		//profile routes
		go privateRoutes.GET("/profile", handlers.GetProfile)
		go privateRoutes.POST("/profile", handlers.SetProfileData)
		go privateRoutes.PUT("/setAge", handlers.SetAge)
		go privateRoutes.PUT("/setEducation", handlers.AddEducation)
		go privateRoutes.PUT("/setPerk", handlers.AddPerk)
		go privateRoutes.PUT("/setDescription", handlers.AddSelfRepresentation)
		go privateRoutes.PUT("/togglePrivacy", handlers.TogglePrivacy)

		go privateRoutes.GET("/messages", handlers.GetMyMessages)
		go privateRoutes.POST("/message", handlers.SendMessage)
	}

	r.Run(c.Port)
}
