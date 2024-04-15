package database

import (
	"fmt"
	"ovo-server/internal/config"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbInst struct {
	connection *gorm.DB
	config     *gorm.Config
}

var db dbInst = dbInst{}

// var db *gorm.DB
// var dbConfig = &gorm.Config{}

func Init() {
	db.config = &gorm.Config{}
	var err error
	switch config.Variables.DatabaseType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Variables.DatabaseUsername,
			config.Variables.DatabasePassword,
			config.Variables.DatabaseHost,
			config.Variables.DatabaseName)

		db.connection, err = gorm.Open(mysql.Open(dsn), db.config)
	default:
		db.connection, err = gorm.Open(sqlite.Open("test.db"), db.config)
	}

	if err != nil {
		panic("Failed to connect to database!")
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
