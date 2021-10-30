package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/johnythelittle/goupdateyourself/configs"
	models "github.com/johnythelittle/goupdateyourself/models/user"
	"github.com/johnythelittle/goupdateyourself/mongoutil"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var user = mongoutil.DB("user")
var config, _ = configs.LoadConfig("../")

func CreateUser(c *gin.Context) {
	var profile = mongoutil.DB("profile")
	var jwtSecret string = config.Secret

	var HashedPassword []byte
	var result models.User
	var credentials models.User

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.Status(500)
	}

	errEmail := user.FindOne(context.TODO(), bson.D{{"email", credentials.Email}}).Decode(&result)
	errUserURL := user.FindOne(context.TODO(), bson.D{{"url_name", credentials.URLName}}).Decode(&result)
	if errEmail != nil {
		if errUserURL != nil {
			HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(credentials.Password), 10)
		} else {
			c.JSON(http.StatusNonAuthoritativeInfo, "user with such url exists already")
			return
		}
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

	ch := make(chan *mongo.InsertOneResult)
	go func(ch chan *mongo.InsertOneResult) {
		i, _ := user.InsertOne(context.TODO(), bson.D{{"username", credentials.Name}, {"password", HashedPassword}, {"email", credentials.Email}, {"url_name", credentials.URLName}})
		ch <- i
	}(ch)
	go func(ch chan *mongo.InsertOneResult) {
		id := <-ch
		profile.InsertOne(context.TODO(), bson.M{"user_id": id.InsertedID.(primitive.ObjectID).Hex(), "is_private": true, "age": 0, "education": []string{}, "perks": []string{}, "description": bson.D{{"text", ""}}, "pronounce": ""})
	}(ch)
}

func Login(c *gin.Context) {
	fmt.Println("CHIPPING IN")
	var credentials models.User
	var result models.User
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
	var userInfo models.User
	var stringToken string = header.Get("authorization")
	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("ERROR ERROR")
		}
		return []byte(config.Secret), nil
	})
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(405, err)
	}
	if token.Valid {
		mapstructure.Decode(token.Claims, &userInfo)
		err := user.FindOne(context.TODO(), bson.D{{"email", userInfo.Email}}).Decode(&userInfo)
		if err != nil {
			fmt.Println("ERROR CANNOT FIND USER", err)
			c.AbortWithError(400, err)
		}
		c.Set("username", userInfo.Name)
		c.Set("email", userInfo.Email)
		c.Set("id", userInfo.ID)
		c.Set("userURL", userInfo.URLName)
	} else {
		c.AbortWithStatusJSON(405, "UNAUTHORIZED")
	}

}

func GetUser(c *gin.Context) {
	//i wrote these lines because of fundamental issue of CORS framework which behaves like a peace of junk wnen i redirect route to group
	//c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500/")
	//c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	//c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	//c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")

	id := c.Query("id")
	var usr models.User
	userIDFormatted, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, err)
	}
	user.FindOne(context.TODO(), bson.D{{"_id", userIDFormatted}}).Decode(&usr)

	usr.Password = ""
	if len(usr.ID) == 0 {
		c.JSON(200, gin.H{"data": false})
	} else {
		c.JSON(200, gin.H{"data": usr})
	}
}

func GetUrl(c *gin.Context) {
	url := c.Query("url")
	var usr models.User
	user.FindOne(context.TODO(), bson.D{{"url_name", url}}).Decode(&usr)
	var userID string = usr.ID
	if userID == "" {
		c.JSON(200, gin.H{"result": false})
	} else {
		c.JSON(200, gin.H{"result": true})
	}
}

func GetMe(c *gin.Context) {
	fmt.Println(">>>>")
	username, _ := c.Get("username")
	email, _ := c.Get("email")
	id, _ := c.Get("id")
	userUrl, _ := c.Get("userURL")
	c.JSON(http.StatusOK, gin.H{"username": username, "email": email, "id": id, "userUrl": userUrl})
}
