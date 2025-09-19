package interfaces

import (
	"backend/routes/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserControllerInterfaces(r *gin.Engine, db *gorm.DB) {
	userRouter := r.Group("/user")

	userRouter.GET("/", func(c *gin.Context) { controllers.FetchUsers(c, db) })
}
