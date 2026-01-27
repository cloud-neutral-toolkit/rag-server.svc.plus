package api

import (
	"gorm.io/gorm"

	"rag-server/internal/model"
	"rag-server/internal/service"
)

// ConfigureServiceDB configures the internal service database connection.
func ConfigureServiceDB(db *gorm.DB) {
	if db != nil {
		_ = db.AutoMigrate(&model.Node{})
	}
	service.SetDB(db)
}
