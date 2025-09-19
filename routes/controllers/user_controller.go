package controllers

import (
	"backend/db_utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchUsers(c *gin.Context, db *gorm.DB) {
	var users []db_utils.User
	res := db.Find(&users)
	if res.Error != nil || res.RowsAffected == 0 {
		c.Abort()
		c.JSON(404, gin.H{"message": "Users not found"})
	}

	c.JSON(200, gin.H{"data": users})
}
