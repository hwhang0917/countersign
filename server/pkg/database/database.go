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
	// Open the file
	file, err := os.Open(datasetPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read and insert words
	for scanner.Scan() {
		word := scanner.Text()

		// Check if the word already exists
		var existingWord models.Word
		result := db.Where("word = ?", word).First(&existingWord)

		if result.Error == gorm.ErrRecordNotFound {
			// Word doesn't exist, so insert it
			newWord := models.Word{Word: word}
			db.Create(&newWord)
			fmt.Printf("Inserted: %s\n", word)
		} else if result.Error != nil {
			log.Printf("Error checking word '%s': %v\n", word, result.Error)
		} else {
			fmt.Printf("Skipped duplicate: %s\n", word)
		}
	}
	return nil
}
