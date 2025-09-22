package migrations

import (
	"search/src/types"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&types.Item{})
}
