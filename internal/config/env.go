package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type variables struct {
	Basedir          string
	ListeningPort    uint
	DatabaseType     string
	DatabaseHost     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string
	TMDBApiKey       string
}

func Get(key string) string {
	value := GetOptional(key)
	if value == "" {
		log.Fatal("Environment variable not set: " + key + ". Exiting.")
	}

	return value
}

func GetOptional(key string) string {
	value := os.Getenv(key)
	return value
}

var (
	Variables variables
)

func Init() {
	// Trying to load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. Using system environment variables.")
	}
	// Setting up the variables
	Variables.DatabaseType = GetOptional("OVO_DATABASE_TYPE")

	if Variables.DatabaseType == "" {
		log.Println("Database type not set. Defaulting to SQLite. This is not recommended for production use. To change this, set the OVO_DATABASE_TYPE environment variables to connect to a MySQL database.")
	}

	if strings.EqualFold("mysql", Variables.DatabaseType) {
		log.Println("Database type set to MySQL. Using MySQL database.")
		Variables.DatabaseHost = Get("OVO_DATABASE_HOST")
		Variables.DatabaseUsername = Get("OVO_DATABASE_USERNAME")
		Variables.DatabasePassword = GetOptional("OVO_DATABASE_PASSWORD")
		Variables.DatabaseName = Get("OVO_DATABASE_NAME")
	}

	if strings.EqualFold("sqlite", Variables.DatabaseType) {
		log.Println("Database type set to SQLite. Using SQLite database.")
	}

	Variables.TMDBApiKey = Get("OVO_TMDB_API_KEY")
	Variables.Basedir = strings.ToLower(GetOptional("OVO_BASEDIR"))
	if GetOptional("OVO_LISTENING_PORT") == "" {
		log.Println("Listening port not set. Using default port 8080.")
		Variables.ListeningPort = 8080
	} else {
		port, err := strconv.Atoi(Get("OVO_LISTENING_PORT"))
		if err != nil {
			log.Fatal("Invalid port number. Using default port 8080.")
			Variables.ListeningPort = 8080
		} else {
			Variables.ListeningPort = uint(port)

		}
	}
}
