package controllers

import (
	"backend/services"
	"net/http"
	"strconv"
	"time"

	"backend/db_utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Yeni session oluştur
func CreateNewSession(c *gin.Context, db *gorm.DB) {
	var input struct {
		UserID   uint `json:"user_id"`
		TargetID uint `json:"target_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := db_utils.Session{
		UserID:    input.UserID,
		TargetID:  input.TargetID,
		StartedAt: time.Now().Unix(),
	}

	if err := db.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// Tüm sessionları listele
func GetSessions(c *gin.Context, db *gorm.DB) {
	var sessions []db_utils.Session
	db.Preload("User").Preload("Target").Find(&sessions)
	c.JSON(http.StatusOK, sessions)
}

// Tek session detayını getir
func GetOneSession(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	var session db_utils.Session
	if err := db.Preload("User").Preload("Target").First(&session, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

// Session update (örn. sadece metadata)
func UpdateSession(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	var session db_utils.Session
	if err := db.First(&session, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	var input struct {
		UserID   uint `json:"user_id"`
		TargetID uint `json:"target_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session.UserID = input.UserID
	session.TargetID = input.TargetID
	db.Save(&session)

	c.JSON(http.StatusOK, session)
}

// Session bitirme
func FinishSession(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	var session db_utils.Session
	if err := db.First(&session, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	session.EndedAt = time.Now().Unix()
	db.Save(&session)

	c.JSON(http.StatusOK, session)
}

// Kullanıcının tüm sessionları
func ListSessionsByUserId(c *gin.Context, db *gorm.DB) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	var sessions []db_utils.Session
	db.Where("user_id = ?", userId).Preload("Target").Find(&sessions)
	c.JSON(http.StatusOK, sessions)
}

// Terminal event kaydet
func SaveSessionData(c *gin.Context, db *gorm.DB) {
	sessionId, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		Type string `json:"type"` // stdin / stdout
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := db_utils.TerminalEvent{
		SessionID: uint(sessionId),
		Type:      input.Type,
		Data:      input.Data,
		Timestamp: time.Now().Unix(),
	}

	if err := db.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// Session eventlerini listele
func GetSessionEvents(c *gin.Context, db *gorm.DB) {
	sessionId, _ := strconv.Atoi(c.Param("id"))
	var events []db_utils.TerminalEvent
	db.Where("session_id = ?", sessionId).Order("timestamp asc").Find(&events)
	c.JSON(http.StatusOK, events)
}

// Son eventleri canlı izleme için getir (örn. son 10)
func GetLatestSessionsEvents(c *gin.Context, db *gorm.DB) {
	sessionId, _ := strconv.Atoi(c.Param("id"))
	var events []db_utils.TerminalEvent
	db.Where("session_id = ?", sessionId).Order("timestamp desc").Limit(10).Find(&events)
	c.JSON(http.StatusOK, events)
}

func GetLive(c *gin.Context, db *gorm.DB) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	var session db_utils.Session
	// Preload ile Target bilgisini de getir
	if err := db.Preload("Target").First(&session, sessionID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Session not found"})
		return
	}

	sshConf := services.SSHConfig{
		User:     session.Target.Username,
		Password: session.Target.Password,
		Host:     session.Target.Host,
		Port:     strconv.Itoa(session.Target.Port),
	}

	services.WSHandler(c, db, sshConf, uint(sessionID))
}
