package config

import (
	"log"
	"os"
	"strconv"

	"github.com/hwhang0917/countersign/internal/constants"
	"github.com/joho/godotenv"
)

func validateConfig() {
	if os.Getenv(constants.API_KEY) == "" {
		log.Fatal("API_KEY is not set in .env file")
	}
}

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env file: ", err)
	}
	validateConfig()
}

func GetPort() string {
	port := os.Getenv(constants.PORT)
	if port == "" {
		port = constants.DEFAULT_PORT
	}
	return port
}

func GetAPIKey() string {
	return os.Getenv(constants.API_KEY)
}

func GetDBFilename() string {
	dbPath := os.Getenv(constants.DB_FILENAME)
	if dbPath == "" {
		dbPath = constants.DEFAULT_DB_FILENAME
	}
	return dbPath
}

func GetInterval() int {
	intervalString := os.Getenv(constants.INTERVAL)
	if intervalString == "" {
		intervalString = constants.DEFAULT_INTERVAL
	}
	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		log.Fatal("Failed to convert interval to integer: ", err)
	}
	return interval
}

func GetDatasetPath() string {
	datasetPath := os.Getenv(constants.DATASET_PATH)
	if datasetPath == "" {
		datasetPath = constants.DEFAULT_DATASET_PATH
	}
	return datasetPath
}
