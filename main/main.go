package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnythelittle/goupdateyourself/configs"
	"github.com/johnythelittle/goupdateyourself/handlers"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	wsserver "github.com/johnythelittle/goupdateyourself/sockets"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}

}
func main() {
	c, err := configs.LoadConfig("../")
	if err != nil {
		log.Panic("BAD CONFIG")
	}
	r := gin.Default()
	r.Use(Options)
	r.Use(CORSMiddleware())
	mongoutil.DB("")

	go wsserver.Manager.Start()
	r.POST("public/login", handlers.Login)
	r.POST("public/register/", handlers.CreateUser)
	go r.GET("content/public_profiles", handlers.GetPublicUsers)
	go r.GET("/ws/", wsserver.ChatHandler)

	go r.GET("public/user_id/", handlers.GetUser)
	go r.GET("get_url/", handlers.GetUrl)
	//publicRoutes.GET("/user/:user_url")
	go r.Use(handlers.CheckUser).GET("me", handlers.GetMe)                          //✅
	go r.Use(handlers.CheckUser).POST("makeBranch", handlers.CreateBranch)          //✅
	go r.Use(handlers.CheckUser).GET("getBranches/", handlers.GetAllBranchesOfUser) //✅
	go r.Use(handlers.CheckUser).PUT("renameBranch", handlers.RenameBranch)
	go r.Use(handlers.CheckUser).PUT("appendBook", handlers.AppendNewElementToBooks)   //✅
	go r.Use(handlers.CheckUser).PUT("modifyBook", handlers.UpdateBookStage)           //✅
	go r.Use(handlers.CheckUser).DELETE("deleteBook", handlers.DeleteElementFromBooks) //✅
	go r.Use(handlers.CheckUser).PUT("appendVideoCourse", handlers.AddVideoCourse)     //✅
	go r.Use(handlers.CheckUser).DELETE("deleteVideoCourse", handlers.DeleteVideoCourse)
	go r.Use(handlers.CheckUser).PUT("appendArticle/", handlers.AddArticle)       //✅
	go r.Use(handlers.CheckUser).DELETE("deleteArticle/", handlers.DeleteArticle) //✅

	//profile routes
	go r.Use(handlers.CheckUser).GET("getMyProfile/", handlers.GetMyProfile) //✅
	go r.Use(handlers.CheckUser).GET("profile/", handlers.GetProfile)        //✅
	go r.Use(handlers.CheckUser).POST("profile", handlers.SetProfileData)    //✅
	go r.Use(handlers.CheckUser).PUT("setAge", handlers.SetAge)              //✅
	go r.Use(handlers.CheckUser).PUT("setEducation", handlers.AddEducation)  //✅
	go r.Use(handlers.CheckUser).DELETE("deleteEducation", handlers.DeleteEducation)
	go r.Use(handlers.CheckUser).PUT("setPerk", handlers.AddPerk) //✅
	go r.Use(handlers.CheckUser).DELETE("deletePerk", handlers.DeletePerk)
	go r.Use(handlers.CheckUser).PUT("setDescription", handlers.AddSelfRepresentation) //✅
	go r.Use(handlers.CheckUser).PUT("togglePrivacy", handlers.TogglePrivacy)          //✅
	go r.Use(handlers.CheckUser).PUT("setPronounce", handlers.SetPronounce)

	//blacklist
	r.POST("/addToBlackList", handlers.AddToBlackList)
	r.DELETE("/deleteFromBlackList", handlers.RemoveFromBlackList)

	//messages

	//r.Use(handlers.CheckUser).GET("/sendMessage", wsserver.ChatHandler)

	r.Run(c.Port)
}
