package database

import (
	"YoutrackGChatBot/settings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDb() (db *gorm.DB, err error) {
	settings, err := settings.GetSettings()

	if err != nil {
		return
	}

	db, err = gorm.Open(sqlite.Open(settings.DATABASE_PATH), &gorm.Config{})

	return
}
