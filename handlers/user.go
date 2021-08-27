package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/johnythelittle/goupdateyourself/configs"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `bson:"_id"`
	Name     string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

var user = mongoutil.DB("user")
var config, _ = configs.LoadConfig("../")

func CreateUser(c *gin.Context) {

	var jwtSecret string = config.Secret

	var HashedPassword []byte
	var result User
	var credentials User

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(500, gin.H{"smthng wrong": "with json parsing", "err": err})
	}

	err := user.FindOne(context.TODO(), bson.D{{"email", credentials.Email}}).Decode(&result)

	if err != nil {
		HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(credentials.Password), 10)
	} else {
		c.JSON(http.StatusNonAuthoritativeInfo, "user with such email exists already")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": credentials.Name,
		"email":    credentials.Email,
	})

	tokenString, _ := token.SignedString([]byte(jwtSecret))

	go func() {
		c.JSON(200, gin.H{
			"token": tokenString,
		})
	}()
	go func() {
		user.InsertOne(context.TODO(), bson.D{{"username", credentials.Name}, {"password", HashedPassword}, {"email", credentials.Email}})
	}()

}

func Login(c *gin.Context) {
	var credentials User
	var result User
	var jwtSecret string = config.Secret
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(500, gin.H{"something wrong": err})
		return
	}
	err := user.FindOne(context.TODO(), bson.D{{"email", credentials.Email}}).Decode(&result)

	if err != nil {
		c.JSON(400, gin.H{"no such user": err})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(credentials.Password)); err != nil {
		log.Println("WRONG HASH")
	} else {
		log.Println("PASSWORD ACCEPTED")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": credentials.Name,
			"email":    credentials.Email,
		})
		stringToken, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{
			"token": stringToken,
		})
	}

}
func CheckUser(c *gin.Context) {
	var header http.Header = c.Request.Header
	var userInfo User
	var stringToken string = header.Get("Authorization")
	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("ERROR ERROR")
		}
		return []byte(config.Secret), nil
	})
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		c.Abort()
	}
	if token.Valid {
		mapstructure.Decode(token.Claims, &userInfo)
		err := user.FindOne(context.TODO(), bson.D{{"email", userInfo.Email}}).Decode(&userInfo)
		if err != nil {
			fmt.Println("error during looking there", err)
			c.AbortWithError(400, err)
		}
		c.Set("username", userInfo.Name)
		c.Set("email", userInfo.Email)
		c.Set("id", userInfo.ID)
	} else {
		c.Abort()
	}

}
