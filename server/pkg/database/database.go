package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hwhang0917/countersign/internal/models"
	"github.com/hwhang0917/countersign/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory: ", err)
	}

	dbPath := filepath.Join(dir, config.GetDBFilename())
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	DB = db
	log.Println("✅ Database connection established")
}

func GetDB() *gorm.DB {
	return DB
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database connection: ", err)
	}
	err = sqlDB.Close()
	if err != nil {
		log.Fatal("Failed to close database connection: ", err)
	}
	log.Println("✅ Database connection closed")
}

func MigrateModels(db *gorm.DB) error {
	log.Println("Migrating models...")
	err := db.AutoMigrate(
		&models.Word{},
	)
	if err != nil {
		log.Println("Failed to migrate models: ", err)
		return err
	}
	log.Println("✅ Models migrated")
	return nil
}
