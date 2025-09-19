package interfaces

import (
	"backend/routes/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SessionControllerInterfaces(r *gin.Engine, db *gorm.DB) {
	sessionRouter := r.Group("/session")

	// CRUD
	sessionRouter.POST("/", func(c *gin.Context) { controllers.CreateNewSession(c, db) })      // Create new session
	sessionRouter.GET("/", func(c *gin.Context) { controllers.GetSessions(c, db) })            // List all sessions
	sessionRouter.GET("/:id", func(c *gin.Context) { controllers.GetOneSession(c, db) })       // Get only one session
	sessionRouter.PUT("/:id", func(c *gin.Context) { controllers.UpdateSession(c, db) })       // Update a session
	sessionRouter.PATCH("/:id/end", func(c *gin.Context) { controllers.FinishSession(c, db) }) // Finish session

	// User related
	sessionRouter.GET("/user/:userId", func(c *gin.Context) { controllers.ListSessionsByUserId(c, db) }) // List sessions by user id

	// Session events (more specific routes first!)
	sessionRouter.POST("/:id/events", func(c *gin.Context) { controllers.SaveSessionData(c, db) }) // Save stdin/stdout data
	sessionRouter.GET("/:id/events", func(c *gin.Context) { controllers.GetSessionEvents(c, db) }) // Get session events
	sessionRouter.GET("/:id/live", func(c *gin.Context) { controllers.GetLive(c, db) })
	sessionRouter.GET("/:id/events/latest", func(c *gin.Context) { controllers.GetLatestSessionsEvents(c, db) }) // Get latest session event for live watching
}
