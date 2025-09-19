package controllers

import (
	"backend/db_utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTargetRequest struct {
	Name     string `json:"name" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateTargetRequest struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     *int   `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateTarget(c *gin.Context, db *gorm.DB) {
	var req CreateTargetRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Println(err)
		c.Abort()
		c.JSON(400, gin.H{"message": "invalid request"})
		return
	}

	newTarget := db_utils.Target{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}

	db.Create(&newTarget)

	c.JSON(201, gin.H{"message": "target created"})
}

func FetchTargets(c *gin.Context, db *gorm.DB) {
	var targets []db_utils.Target
	res := db.Find(&targets)
	if res.Error != nil || res.RowsAffected == 0 {
		c.Abort()
		c.JSON(404, gin.H{"message": "targets not found"})
		return
	}

	c.JSON(200, gin.H{"data": targets})
}

func FetchOneTarget(c *gin.Context, db *gorm.DB) {
	targetId := c.Param("targetId")
	var target db_utils.Target
	res := db.Where("id = ?", targetId).First(&target)
	if res.Error != nil || res.RowsAffected == 0 {
		c.Abort()
		c.JSON(404, gin.H{"message": "target not found"})
		return
	}

	c.JSON(200, gin.H{"data": target})
}

func UpdateTarget(c *gin.Context, db *gorm.DB) {
	var req UpdateTargetRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Println(err)
		c.Abort()
		c.JSON(400, gin.H{"message": "invalid request"})
		return
	}

	targetId := c.Param("targetId")

	var currentTarget db_utils.Target
	res := db.Where("id = ?", targetId).First(&currentTarget)
	if res.Error != nil || res.RowsAffected == 0 {
		c.Abort()
		c.JSON(404, gin.H{"message": "target not found"})
		return
	}

	if req.Name != "" {
		currentTarget.Name = req.Name
	}
	if req.Host != "" {
		currentTarget.Host = req.Host
	}
	if req.Port != nil || *req.Port != 0 {
		currentTarget.Port = *req.Port
	}
	if req.Username != "" {
		currentTarget.Username = req.Username
	}
	if req.Password != "" {
		currentTarget.Password = req.Password
	}
	db.Save(&currentTarget)

	c.JSON(200, gin.H{"data": currentTarget})
}

func DeleteTarget(c *gin.Context, db *gorm.DB) {
	targetId := c.Param("targetId")

	var target db_utils.Target
	res := db.Where("id = ?", targetId).First(&target)
	if res.Error != nil || res.RowsAffected == 0 {
		c.Abort()
		c.JSON(404, gin.H{"message": "target not found"})
		return
	}

	db.Delete(&target)

	c.JSON(200, gin.H{"message": "Deleted successfully"})
}
