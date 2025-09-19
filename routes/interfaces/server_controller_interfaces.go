package interfaces

import (
	"backend/routes/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ServerControllerInterfaces(r *gin.Engine, db *gorm.DB) {
	serverRouter := r.Group("/server")

	serverRouter.POST("/", func(c *gin.Context) { controllers.CreateTarget(c, db) })            // Create a new target
	serverRouter.GET("/", func(c *gin.Context) { controllers.FetchTargets(c, db) })             // Get Targets
	serverRouter.GET("/:targetId", func(c *gin.Context) { controllers.FetchOneTarget(c, db) })  // Get only one target
	serverRouter.PUT("/:targetId", func(c *gin.Context) { controllers.UpdateTarget(c, db) })    // Update target
	serverRouter.DELETE("/:targetId", func(c *gin.Context) { controllers.DeleteTarget(c, db) }) // Delete Target
}
