package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnythelittle/goupdateyourself/configs"
	"github.com/johnythelittle/goupdateyourself/handlers"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
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

	r.GET("content/public_profiles", handlers.GetPublicUsers)

	r.POST("public/login", handlers.Login)
	r.POST("public/register/", handlers.CreateUser)
	r.GET("public/user_id/", handlers.GetUser)
	r.GET("get_url/", handlers.GetUrl)
	//publicRoutes.GET("/user/:user_url")

	r.Use(handlers.CheckUser).GET("me/", handlers.GetMe)                //✅
	r.Use(handlers.CheckUser).POST("makeBranch", handlers.CreateBranch) //✅
	r.Use(handlers.CheckUser).GET("getBranches/", handlers.GetAllBranchesOfUser)
	r.Use(handlers.CheckUser).PUT("renameBranch", handlers.RenameBranch)
	r.Use(handlers.CheckUser).PUT("appendBook", handlers.AppendNewElementToBooks)
	r.Use(handlers.CheckUser).PUT("modifyBook", handlers.UpdateBookStage)
	r.Use(handlers.CheckUser).PUT("deleteBook", handlers.DeleteElementFromBooks)
	r.Use(handlers.CheckUser).PUT("appendVideoCourse", handlers.AddVideoCourse)

	//profile routes
	r.Use(handlers.CheckUser).GET("getMyProfile/", handlers.GetMyProfile) //✅
	r.Use(handlers.CheckUser).GET("profile", handlers.GetProfile)         //✅
	r.Use(handlers.CheckUser).POST("profile", handlers.SetProfileData)    //✅
	r.Use(handlers.CheckUser).PUT("setAge", handlers.SetAge)              //✅
	r.Use(handlers.CheckUser).PUT("setEducation", handlers.AddEducation)  //✅
	r.Use(handlers.CheckUser).DELETE("deleteEducation", handlers.DeleteEducation)
	r.Use(handlers.CheckUser).PUT("setPerk", handlers.AddPerk) //✅
	r.Use(handlers.CheckUser).DELETE("deletePerk", handlers.DeletePerk)
	r.Use(handlers.CheckUser).PUT("setDescription", handlers.AddSelfRepresentation) //✅
	r.Use(handlers.CheckUser).PUT("togglePrivacy", handlers.TogglePrivacy)          //✅
	r.Use(handlers.CheckUser).PUT("setPronounce", handlers.SetPronounce)

	//blacklist
	r.POST("/addToBlackList", handlers.AddToBlackList)
	r.DELETE("/deleteFromBlackList", handlers.RemoveFromBlackList)

	//messages
	r.POST("/addDialogue", handlers.AddDialogue)
	r.POST("/sendMessage", handlers.SendMessage)
	r.GET("/getMyMessages", handlers.GetMyDialogues)

	r.Run(c.Port)
}
