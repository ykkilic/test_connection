package main

import (
	"backend/db_utils"
	"backend/routes/interfaces"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Server Error: ", err)
		}
	}()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // izin verilen frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db := db_utils.Connect()
	db_utils.CreateTables(db)

	interfaces.ServerControllerInterfaces(router, db)
	interfaces.SessionControllerInterfaces(router, db)
	interfaces.UserControllerInterfaces(router, db)

	err := router.Run("localhost:8090")
	if err != nil {
		panic(err)
	}
}
