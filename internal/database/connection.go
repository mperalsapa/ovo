package database

import (
	"fmt"
	"log"
	"os"
	"ovo-server/internal/config"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbInst struct {
	connection *gorm.DB
	config     *gorm.Config
}

var db dbInst = dbInst{}

func Init() {
	db.config = &gorm.Config{}
	var err error

	// Copied from GORM Docs
	db.config.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	tries := 0
	delay := 6

	for tries < 5 {
		if tries > 0 {
			log.Println("Retrying database connection in ", delay, " seconds")
			time.Sleep(time.Duration(delay) * time.Second)
		}

		switch config.Variables.DatabaseType {
		case "mysql":
			dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				config.Variables.DatabaseUsername,
				config.Variables.DatabasePassword,
				config.Variables.DatabaseHost,
				config.Variables.DatabaseName)

			db.connection, err = gorm.Open(mysql.Open(dsn), db.config)
		default:
			db.connection, err = gorm.Open(sqlite.Open("ovo.db"), db.config)
		}
		if err != nil {
			log.Println("Failed to connect to database: ", err.Error())
			tries++
			if tries == 5 {
				log.Println("Failed to connect to database after 5 tries. Exiting.")
				os.Exit(1)
			}
		} else {
			break
		}
	}

}

func GetDB() *gorm.DB {
	return db.connection
}

func Create(model interface{}) {
	db.connection.Create(model)
}

func Migrate(models ...interface{}) {
	db.connection.AutoMigrate(models...)
}
