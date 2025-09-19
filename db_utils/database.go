package db_utils

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Veritabanı Bağlantısında bir problem oluştu: ", err)
			return
		}
	}()

	dsn := "host=localhost user=postgres password=Bjk1903 dbname=a_db port=5432 sslmode=disable TimeZone=Europe/Istanbul"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func CreateTables(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Target{}, &Session{}, &TerminalEvent{})
}
