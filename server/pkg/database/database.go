package database

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hwhang0917/countersign/internal/models"
	"github.com/hwhang0917/countersign/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func PopulateWords(db *gorm.DB) error {
	datasetPath := config.GetDatasetPath()
	log.Printf("Populating words from %s...", datasetPath)

	currentLogger := db.Logger
	db.Logger = db.Logger.LogMode(logger.Silent)
	defer func() {
		db.Logger = currentLogger
	}()

	// Open the file
	file, err := os.Open(datasetPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		insertWord(db, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading file: [%s] due to %v", datasetPath, err)
	}
	return nil
}

func insertWord(db *gorm.DB, word string) error {
	var existingWord models.Word
	result := db.Where("word = ?", word).First(&existingWord)
	if result.Error == gorm.ErrRecordNotFound {
		newWord := models.Word{Word: word}
		if err := db.Create(&newWord).Error; err != nil {
			return fmt.Errorf("Failed to insert word [%s]: %v", word, err)
		}
		// Print progress in same line
		fmt.Printf("\rPopulating words... [ID: %d, Word: %s]", newWord.ID, newWord.Word)
	} else if result.Error != nil {
		return fmt.Errorf("Failed to query word [%s]: %v", word, result.Error)
	} else {
		return nil
	}
	return nil
}
