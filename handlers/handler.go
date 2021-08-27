package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	fmt.Println("yeah, thats how it exactly works")
	c.JSON(200, gin.H{"you see?": "its really simple"})
}
