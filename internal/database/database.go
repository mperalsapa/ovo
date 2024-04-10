package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	connection *gorm.DB
	config     *gorm.Config
}

var db DB = DB{}

// var db *gorm.DB
// var dbConfig = &gorm.Config{}

func Init() {
	db.config = &gorm.Config{}
	var err error
	switch "mysql" {
	case "mysql":
		dsn := "root:@tcp(127.0.0.1:3306)/ovo-dev?charset=utf8mb4&parseTime=True&loc=Local"
		db.connection, err = gorm.Open(mysql.Open(dsn), db.config)
	case "sqlite":
		db.connection, err = gorm.Open(sqlite.Open("test.db"), db.config)
	}

	if err != nil {
		panic("Failed to connect to database!")
	}

	Migrate()
}
