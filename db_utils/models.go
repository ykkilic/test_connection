package db_utils

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt int64
}

type Target struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Host     string `gorm:"not null"` // IP veya hostname
	Port     int    `gorm:"default:22"`
	Username string `gorm:"not null"` // hedef kullanıcı
	Password string // opsiyonel, vault ile değiştirilebilir
}

type Session struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User
	TargetID  uint
	Target    Target
	StartedAt int64
	EndedAt   int64
}

type TerminalEvent struct {
	ID        uint `gorm:"primaryKey"`
	SessionID uint
	Session   Session
	Type      string `gorm:"not null"` // "stdin" veya "stdout"
	Data      string `gorm:"type:text;not null"`
	Timestamp int64  `gorm:"autoCreateTime"`
}
